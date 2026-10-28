package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func postJSON(h http.HandlerFunc, body any) *httptest.ResponseRecorder {
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(b))
	rec := httptest.NewRecorder()
	h(rec, req)
	return rec
}

func installForTest(t *testing.T, h *AppHandler, insalesID string) {
	t.Helper()
	req := httptest.NewRequest("GET", "/install?shop=shop.myinsales.ru&token=tok&insales_id="+insalesID, nil)
	rec := httptest.NewRecorder()
	h.Install(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("install: %d", rec.Code)
	}
}

func TestRegister_VerifyBeforeUse(t *testing.T) {
	h, repo := setupAuthHandler()
	if _, err := repo.Create("existing", "old@test.com", "", "", "", "hash"); err != nil {
		t.Fatal(err)
	}

	rec := postJSON(h.Register, RegisterRequest{Nickname: "newbie", Email: "New@TEST.com", Password: "RegisterPass1"})
	if rec.Code != http.StatusCreated {
		t.Fatalf("register: %d %s", rec.Code, rec.Body.String())
	}
	var resp AuthResponse
	_ = json.Unmarshal(rec.Body.Bytes(), &resp)
	if !resp.RequireEmail {
		t.Error("register must require email verification")
	}

	u := repo.users[resp.UserID]
	if u == nil {
		t.Fatal("user not created")
	}
	if u.Email != "new@test.com" {
		t.Errorf("email = %q, want lowercase", u.Email)
	}
	if u.EmailVerified {
		t.Error("new registration must be unverified")
	}

	// У «мёртвого» аккаунта нет известного пароля — вход невозможен до верификации.
	rec = postJSON(h.Login, LoginRequest{Login: "newbie", Password: "RegisterPass1"})
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("login before verify: %d %s", rec.Code, rec.Body.String())
	}

	// Гейт: даже зная пароль (атакующий задал его при регистрации), вход закрыт.
	attHash, _ := bcrypt.GenerateFromPassword([]byte("AttackerPass1"), bcrypt.DefaultCost)
	if err := repo.UpdatePassword(resp.UserID, string(attHash)); err != nil {
		t.Fatal(err)
	}
	rec = postJSON(h.Login, LoginRequest{Login: "newbie", Password: "AttackerPass1"})
	if rec.Code != http.StatusForbidden {
		t.Fatalf("gate must block unverified login: %d %s", rec.Code, rec.Body.String())
	}

	rec = postJSON(h.VerifyEmail, VerifyEmailRequest{UserID: resp.UserID, Code: "000000", Password: "OwnerPass1"})
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("bad code: %d %s", rec.Code, rec.Body.String())
	}

	// без пароля — отказ (пароль ставит только подтвердивший)
	code := h.codes.m["register:"+resp.UserID].code
	rec = postJSON(h.VerifyEmail, VerifyEmailRequest{UserID: resp.UserID, Code: code})
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("verify without password: %d %s", rec.Code, rec.Body.String())
	}

	code = h.codes.m["register:"+resp.UserID].code
	rec = postJSON(h.VerifyEmail, VerifyEmailRequest{UserID: resp.UserID, Code: code, Password: "OwnerPass1"})
	if rec.Code != http.StatusOK {
		t.Fatalf("verify: %d %s", rec.Code, rec.Body.String())
	}
	if !repo.users[resp.UserID].EmailVerified {
		t.Error("email must be verified after code")
	}

	// пароль, выставленный при регистрации, не действует — работает пароль владельца
	rec = postJSON(h.Login, LoginRequest{Login: "newbie", Password: "RegisterPass1"})
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("register password must be dead: %d %s", rec.Code, rec.Body.String())
	}
	rec = postJSON(h.Login, LoginRequest{Login: "newbie", Password: "OwnerPass1"})
	if rec.Code != http.StatusOK {
		t.Fatalf("login after verify: %d %s", rec.Code, rec.Body.String())
	}
}

func TestLinkCode_VerifiedAccount_NoSeize(t *testing.T) {
	h, installs, users, _, _ := newTestAppHandler(t)
	installForTest(t, h, "999")

	hash, err := bcrypt.GenerateFromPassword([]byte("OldPass123"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}
	oldHash := string(hash)
	user, err := users.Create("victim", "victim@test.com", "", "", "", oldHash)
	if err != nil {
		t.Fatal(err)
	}
	if err := installs.SetPendingLink("999", user.ID); err != nil {
		t.Fatal(err)
	}

	rec := postJSON(h.LinkCodeSend, map[string]string{"insales_id": "999"})
	if rec.Code != http.StatusOK {
		t.Fatalf("send: %d %s", rec.Code, rec.Body.String())
	}
	var sendResp map[string]any
	_ = json.Unmarshal(rec.Body.Bytes(), &sendResp)
	if sendResp["email"] != "vi***@test.com" {
		t.Errorf("email must be masked, got %v", sendResp["email"])
	}

	code := h.codes.m["link:999"].code
	rec = postJSON(h.LinkCodeVerify, map[string]string{"insales_id": "999", "code": code})
	if rec.Code != http.StatusOK {
		t.Fatalf("verify: %d %s", rec.Code, rec.Body.String())
	}
	var verifyResp map[string]any
	_ = json.Unmarshal(rec.Body.Bytes(), &verifyResp)
	if verifyResp["code"] == nil || verifyResp["code"] == "" {
		t.Error("verify must return exchange code")
	}
	if verifyResp["require_password_change"] != false {
		t.Error("verified account must keep its password (no seize)")
	}

	updated := users.users[user.ID]
	if updated.PasswordHash != oldHash {
		t.Error("password of verified account must not change on link")
	}
	if updated.RequirePasswordChange {
		t.Error("verified account must not be forced to change password")
	}
	inst, _ := installs.GetByInsalesID("999")
	if inst.ShopID == "" {
		t.Error("shop must be provisioned on link")
	}
	if inst.UserID != user.ID {
		t.Errorf("install bound to %s, want %s", inst.UserID, user.ID)
	}
}

func TestLinkCode_DormantAccount_RevokeAndForceChange(t *testing.T) {
	h, installs, users, _, _ := newTestAppHandler(t)
	installForTest(t, h, "998")

	hash, _ := bcrypt.GenerateFromPassword([]byte("AttackerPass1"), bcrypt.DefaultCost)
	oldHash := string(hash)
	// «Мёртвый» аккаунт: зарегистрирован без подтверждения — пароль чужой/нерабочий.
	user, _ := users.Create("dormant", "dormant@test.com", "", "", "", oldHash)
	_ = users.UpdateEmailVerified(user.ID, false)
	_ = installs.SetPendingLink("998", user.ID)

	code := func() string {
		rec := postJSON(h.LinkCodeSend, map[string]string{"insales_id": "998"})
		if rec.Code != http.StatusOK {
			t.Fatalf("send: %d", rec.Code)
		}
		return h.codes.m["link:998"].code
	}()

	rec := postJSON(h.LinkCodeVerify, map[string]string{"insales_id": "998", "code": code})
	if rec.Code != http.StatusOK {
		t.Fatalf("verify: %d %s", rec.Code, rec.Body.String())
	}
	var verifyResp map[string]any
	_ = json.Unmarshal(rec.Body.Bytes(), &verifyResp)
	if verifyResp["require_password_change"] != true {
		t.Error("dormant account must force password set")
	}

	updated := users.users[user.ID]
	if updated.PasswordHash == oldHash {
		t.Error("dormant password must be revoked")
	}
	if !updated.RequirePasswordChange {
		t.Error("require_password_change must be set")
	}
	if !updated.EmailVerified {
		t.Error("link code proves inbox — email must become verified")
	}
}

func TestLinkCode_VerifyInvalid(t *testing.T) {
	h, installs, users, _, _ := newTestAppHandler(t)
	installForTest(t, h, "777")

	hash, _ := bcrypt.GenerateFromPassword([]byte("OldPass123"), bcrypt.DefaultCost)
	oldHash := string(hash)
	user, _ := users.Create("u1", "u1@test.com", "", "", "", oldHash)
	_ = installs.SetPendingLink("777", user.ID)

	rec := postJSON(h.LinkCodeVerify, map[string]string{"insales_id": "777", "code": "123456"})
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("bad code: %d %s", rec.Code, rec.Body.String())
	}
	if users.users[user.ID].PasswordHash != oldHash {
		t.Error("password must not change on invalid code")
	}
}

func TestLinkCode_SendCooldown(t *testing.T) {
	h, installs, users, _, _ := newTestAppHandler(t)
	installForTest(t, h, "555")

	hash, _ := bcrypt.GenerateFromPassword([]byte("OldPass123"), bcrypt.DefaultCost)
	user, _ := users.Create("u2", "u2@test.com", "", "", "", string(hash))
	_ = installs.SetPendingLink("555", user.ID)

	if rec := postJSON(h.LinkCodeSend, map[string]string{"insales_id": "555"}); rec.Code != http.StatusOK {
		t.Fatalf("first send: %d", rec.Code)
	}
	rec := postJSON(h.LinkCodeSend, map[string]string{"insales_id": "555"})
	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("cooldown: %d %s", rec.Code, rec.Body.String())
	}
}
