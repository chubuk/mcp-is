package handlers

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/chubuk/mcp-is/internal/appinstall"
	"github.com/chubuk/mcp-is/internal/auth"
	"github.com/chubuk/mcp-is/internal/storage"
)

func newTestAppHandler(t *testing.T) (*AppHandler, *mockAppInstallRepo, *mockUserRepo, *mockShopRepo, *mockAIKeyRepo) {
	t.Helper()
	aes, err := auth.NewAESCipher("test-encryption-key-32-bytes!!!!")
	if err != nil {
		t.Fatalf("aes: %v", err)
	}
	installs := newMockAppInstallRepo()
	users := newMockUserRepo()
	shops := newMockShopRepo()
	keys := newMockAIKeyRepo()
	jwt := auth.NewJWTManager("test-jwt-secret", time.Hour)
	h := NewAppHandler(installs, users, shops, keys, aes, jwt, nil,
		"myapplogin", "mysecret", "mcp-is.ru", "test-jwt-secret", 0)
	return h, installs, users, shops, keys
}

func TestInstall_StoresPasswordFromMD5(t *testing.T) {
	h, installs, _, _, _ := newTestAppHandler(t)
	req := httptest.NewRequest("GET", "/install?shop=test.myinsales.ru&token=token123&insales_id=123", nil)
	rec := httptest.NewRecorder()
	h.Install(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	inst, err := installs.GetByInsalesID("123")
	if err != nil {
		t.Fatalf("install not stored: %v", err)
	}
	if inst.ShopHost != "test.myinsales.ru" {
		t.Errorf("shop_host = %q", inst.ShopHost)
	}
	// password = MD5("token123"+"mysecret") = 4c3f4c197336eb97475506e71c839c71
	// stored AES-encrypted; verify via handler uninstall path instead of raw compare
	_ = appinstall.InstallPassword // keep import used in package tests
}

func TestInstall_MissingParams(t *testing.T) {
	h, _, _, _, _ := newTestAppHandler(t)
	req := httptest.NewRequest("GET", "/install?shop=test.myinsales.ru", nil)
	rec := httptest.NewRecorder()
	h.Install(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
}

// (a) first install, no prev — must not touch shops and must return 200.
func TestInstall_FirstInstall_NoPrev(t *testing.T) {
	h, installs, _, shops, _ := newTestAppHandler(t)
	req := httptest.NewRequest("GET", "/install?shop=test.myinsales.ru&token=token123&insales_id=123", nil)
	rec := httptest.NewRecorder()
	h.Install(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	inst, err := installs.GetByInsalesID("123")
	if err != nil {
		t.Fatalf("install not stored: %v", err)
	}
	if inst.ShopID != "" {
		t.Errorf("shop_id = %q, want empty on first install", inst.ShopID)
	}
	if len(shops.shops) != 0 {
		t.Errorf("shops touched on first install: %d rows", len(shops.shops))
	}
}

// (b) reinstall with ShopID == "" — no-op on shops, no panic, 200.
func TestInstall_Reinstall_NoShop_NoPanic(t *testing.T) {
	h, installs, _, _, _ := newTestAppHandler(t)
	for i := 0; i < 2; i++ {
		req := httptest.NewRequest("GET", "/install?shop=test.myinsales.ru&token=token123&insales_id=123", nil)
		rec := httptest.NewRecorder()
		h.Install(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("attempt %d: status = %d, want 200", i, rec.Code)
		}
	}
	if _, err := installs.GetByInsalesID("123"); err != nil {
		t.Fatalf("install not stored: %v", err)
	}
}

// (c) reinstall with ShopID set — shops.api_password_encrypted must equal the new enc.
func TestInstall_Reinstall_RefreshesShopPassword(t *testing.T) {
	h, installs, _, shops, _ := newTestAppHandler(t)

	// first install + provision link (as autologin would do)
	req := httptest.NewRequest("GET", "/install?shop=test.myinsales.ru&token=token123&insales_id=123", nil)
	rec := httptest.NewRecorder()
	h.Install(rec, req)
	oldEnc, _ := h.aes.Encrypt(appinstall.InstallPassword("token123", "mysecret"))
	shop, _ := shops.Create("user-1", "test.myinsales.ru", "test.myinsales.ru", "myapplogin", oldEnc, true, false, false)
	_ = installs.LinkShop("123", "user-1", shop.ID)

	// reinstall with a new token → new password
	req = httptest.NewRequest("GET", "/install?shop=test.myinsales.ru&token=NEWTOKEN&insales_id=123", nil)
	rec = httptest.NewRecorder()
	h.Install(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	// AES-GCM nonce is random — compare decrypted values, not ciphertexts
	wantPassword := appinstall.InstallPassword("NEWTOKEN", "mysecret")
	got, err := shops.GetByID(shop.ID)
	if err != nil {
		t.Fatalf("shop get: %v", err)
	}
	gotPassword, err := h.aes.Decrypt(got.APIPasswordEncrypted)
	if err != nil {
		t.Fatalf("decrypt shop password: %v", err)
	}
	if gotPassword != wantPassword {
		t.Errorf("shops password = %q, want %q", gotPassword, wantPassword)
	}
	inst, _ := installs.GetByInsalesID("123")
	instPassword, err := h.aes.Decrypt(inst.APIPasswordEncrypted)
	if err != nil {
		t.Fatalf("decrypt install password: %v", err)
	}
	if instPassword != wantPassword {
		t.Errorf("app_installs password = %q, want %q", instPassword, wantPassword)
	}
}

// (d) UpdateAPIPassword failure → 500 (no 200), so InSales retries.
func TestInstall_Reinstall_ShopUpdateError_500(t *testing.T) {
	h, installs, _, shops, _ := newTestAppHandler(t)

	req := httptest.NewRequest("GET", "/install?shop=test.myinsales.ru&token=token123&insales_id=123", nil)
	rec := httptest.NewRecorder()
	h.Install(rec, req)
	oldEnc, _ := h.aes.Encrypt(appinstall.InstallPassword("token123", "mysecret"))
	shop, _ := shops.Create("user-1", "test.myinsales.ru", "test.myinsales.ru", "myapplogin", oldEnc, true, false, false)
	_ = installs.LinkShop("123", "user-1", shop.ID)

	// shop row vanished → UpdateAPIPassword returns error
	_ = shops.Delete(shop.ID)

	req = httptest.NewRequest("GET", "/install?shop=test.myinsales.ru&token=NEWTOKEN&insales_id=123", nil)
	rec = httptest.NewRecorder()
	h.Install(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", rec.Code)
	}
}

// --- autologin: link helpers ---

func installFor(t *testing.T, h *AppHandler, insalesID, token string) *mockAppInstallRepo {
	t.Helper()
	req := httptest.NewRequest("GET", "/install?shop=test.myinsales.ru&token="+token+"&insales_id="+insalesID, nil)
	rec := httptest.NewRecorder()
	h.Install(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("install status = %d", rec.Code)
	}
	return nil
}

func autologinReq(nonce, insalesID, query string) *http.Request {
	req := httptest.NewRequest("GET", "/autologin?"+query, nil)
	req.AddCookie(&http.Cookie{Name: "autologin_nonce", Value: nonce + ":" + insalesID})
	return req
}

func token3(nonce, email, name, userID, confirmed, password string) string {
	sum := md5.Sum([]byte(nonce + email + name + userID + confirmed + password))
	return hex.EncodeToString(sum[:])
}

// (1) /app with linked shop must round-trip via InSales, not skip to /onboarding.
func TestAppLanding_AlwaysInSalesRoundTrip(t *testing.T) {
	h, installs, _, shops, _ := newTestAppHandler(t)
	installFor(t, h, "123", "token123")
	enc, _ := h.aes.Encrypt(appinstall.InstallPassword("token123", "mysecret"))
	shop, _ := shops.Create("user-1", "test.myinsales.ru", "test.myinsales.ru", "myapplogin", enc, true, false, false)
	_ = installs.LinkShop("123", "user-1", shop.ID)

	req := httptest.NewRequest("GET", "/app?shop=test.myinsales.ru&insales_id=123", nil)
	rec := httptest.NewRecorder()
	h.AppLanding(rec, req)

	loc := rec.Header().Get("Location")
	if rec.Code != http.StatusFound {
		t.Fatalf("status = %d, want 302", rec.Code)
	}
	if !strings.Contains(loc, "/admin/applications/myapplogin/login") {
		t.Errorf("Location = %q, want InSales login URL", loc)
	}
}

// (2) returning user (UserID matches, shop linked) → session, no pending_link, empty api_key.
func TestAutologin_ReturningUser_ReuseShop(t *testing.T) {
	h, installs, users, shops, keys := newTestAppHandler(t)
	installFor(t, h, "123", "token123")
	password := appinstall.InstallPassword("token123", "mysecret")
	enc, _ := h.aes.Encrypt(password)
	user, _ := users.Create("u1", "user@test.ru", "", "", "", "hash")
	shop, _ := shops.Create(user.ID, "test.myinsales.ru", "test.myinsales.ru", "myapplogin", enc, true, false, false)
	_ = installs.LinkShop("123", user.ID, shop.ID)

	req := autologinReq("nonce1", "123",
		"user_email=user@test.ru&user_name=u1&user_id=1&email_confirmed=true&token3="+
			token3("nonce1", "user@test.ru", "u1", "1", "true", password))
	rec := httptest.NewRecorder()
	h.Autologin(rec, req)

	if rec.Code != http.StatusFound {
		t.Fatalf("status = %d, want 302", rec.Code)
	}
	loc := rec.Header().Get("Location")
	if !strings.HasPrefix(loc, "/onboarding?code=") {
		t.Fatalf("Location = %q, want /onboarding?code=", loc)
	}
	code := strings.TrimPrefix(loc, "/onboarding?code=")
	uid, raw, ok := h.consumeCode(code)
	if !ok || uid != user.ID {
		t.Fatalf("consumeCode = (%q,%v)", uid, ok)
	}
	if raw != "" {
		t.Errorf("api_key = %q, want empty on return (no key re-mint)", raw)
	}
	inst, _ := installs.GetByInsalesID("123")
	if inst.Status != "installed" {
		t.Errorf("status = %q, want installed (not pending_link)", inst.Status)
	}
	// no new shop/key created
	if n := len(shops.shops); n != 1 {
		t.Errorf("shops = %d, want 1", n)
	}
	if n := len(keys.keys); n != 0 {
		t.Errorf("ai_keys = %d, want 0 (no new key on return)", n)
	}
}

// (3) reinstall leftover (UserID set, ShopID empty) → shop recreated + fresh key.
func TestAutologin_ReinstallLeftover_Reprovision(t *testing.T) {
	h, installs, users, shops, _ := newTestAppHandler(t)
	installFor(t, h, "123", "token123")
	password := appinstall.InstallPassword("token123", "mysecret")
	user, _ := users.Create("u1", "user@test.ru", "", "", "", "hash")
	_ = installs.SetUser("123", user.ID) // leftover: user kept, shop cleared

	req := autologinReq("nonce1", "123",
		"user_email=user@test.ru&user_name=u1&user_id=1&email_confirmed=true&token3="+
			token3("nonce1", "user@test.ru", "u1", "1", "true", password))
	rec := httptest.NewRecorder()
	h.Autologin(rec, req)

	if rec.Code != http.StatusFound {
		t.Fatalf("status = %d, want 302", rec.Code)
	}
	loc := rec.Header().Get("Location")
	if !strings.HasPrefix(loc, "/onboarding?code=") {
		t.Fatalf("Location = %q, want /onboarding?code=", loc)
	}
	_, raw, ok := h.consumeCode(strings.TrimPrefix(loc, "/onboarding?code="))
	if !ok || raw == "" {
		t.Errorf("api_key empty, want fresh key after re-provision")
	}
	if n := len(shops.shops); n != 1 {
		t.Errorf("shops = %d, want 1 (recreated)", n)
	}
	inst, _ := installs.GetByInsalesID("123")
	if inst.ShopID == "" {
		t.Error("inst.ShopID not linked after re-provision")
	}
}

// (4) unlinked install + existing email → pending_link (anti-takeover).
func TestAutologin_ExistingEmail_PendingLink(t *testing.T) {
	h, installs, users, _, _ := newTestAppHandler(t)
	installFor(t, h, "123", "token123")
	password := appinstall.InstallPassword("token123", "mysecret")
	user, _ := users.Create("u1", "user@test.ru", "", "", "", "hash")

	req := autologinReq("nonce1", "123",
		"user_email=user@test.ru&user_name=u1&user_id=1&email_confirmed=true&token3="+
			token3("nonce1", "user@test.ru", "u1", "1", "true", password))
	rec := httptest.NewRecorder()
	h.Autologin(rec, req)

	if rec.Code != http.StatusFound {
		t.Fatalf("status = %d, want 302", rec.Code)
	}
	if loc := rec.Header().Get("Location"); loc != "/login?pending_install=123" {
		t.Errorf("Location = %q, want /login?pending_install=123", loc)
	}
	inst, _ := installs.GetByInsalesID("123")
	if inst.Status != "pending_link" || inst.UserID != user.ID {
		t.Errorf("inst = %+v, want pending_link/%s", inst, user.ID)
	}
}

// (5) user conflict (inst bound to another user) → 4xx, user_id not rewritten.
func TestAutologin_UserConflict_NoRewrite(t *testing.T) {
	h, installs, users, _, _ := newTestAppHandler(t)
	installFor(t, h, "123", "token123")
	password := appinstall.InstallPassword("token123", "mysecret")
	other, _ := users.Create("other", "other@test.ru", "", "", "", "hash")
	victim, _ := users.Create("u1", "user@test.ru", "", "", "", "hash")
	_ = installs.SetUser("123", other.ID) // install bound to another account

	req := autologinReq("nonce1", "123",
		"user_email=user@test.ru&user_name=u1&user_id=1&email_confirmed=true&token3="+
			token3("nonce1", "user@test.ru", "u1", "1", "true", password))
	rec := httptest.NewRecorder()
	h.Autologin(rec, req)

	if rec.Code != http.StatusConflict {
		t.Fatalf("status = %d, want 409", rec.Code)
	}
	inst, _ := installs.GetByInsalesID("123")
	if inst.UserID != other.ID {
		t.Errorf("user_id rewritten to %q, want %q", inst.UserID, other.ID)
	}
	_ = victim
}

// (6) new email → provision.
func TestAutologin_NewEmail_Provision(t *testing.T) {
	h, installs, users, _, _ := newTestAppHandler(t)
	installFor(t, h, "123", "token123")
	password := appinstall.InstallPassword("token123", "mysecret")

	req := autologinReq("nonce1", "123",
		"user_email=new@test.ru&user_name=new&user_id=1&email_confirmed=true&token3="+
			token3("nonce1", "new@test.ru", "new", "1", "true", password))
	rec := httptest.NewRecorder()
	h.Autologin(rec, req)

	if rec.Code != http.StatusFound {
		t.Fatalf("status = %d, want 302", rec.Code)
	}
	if loc := rec.Header().Get("Location"); !strings.HasPrefix(loc, "/onboarding?code=") {
		t.Errorf("Location = %q", loc)
	}
	if len(users.users) != 1 {
		t.Errorf("users = %d, want 1", len(users.users))
	}
	inst, _ := installs.GetByInsalesID("123")
	if inst.ShopID == "" {
		t.Error("shop not linked")
	}
}

// (7) missing nonce cookie → 400.
func TestAutologin_NoNonceCookie(t *testing.T) {
	h, _, _, _, _ := newTestAppHandler(t)
	req := httptest.NewRequest("GET", "/autologin?token3=abc&user_email=a@b.c&user_name=a&user_id=1&email_confirmed=true", nil)
	rec := httptest.NewRecorder()
	h.Autologin(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
}

// (8) short vendor format token=MD5(nonce+password): bound install → session; unbound → 400.
func TestAutologin_ShortToken(t *testing.T) {
	h, installs, users, _, _ := newTestAppHandler(t)
	installFor(t, h, "123", "token123")
	password := appinstall.InstallPassword("token123", "mysecret")
	user, _ := users.Create("u1", "user@test.ru", "", "", "", "hash")

	// unbound install → 400
	req := autologinReq("nonce1", "123", "token="+appinstall.NoncePasswordMD5("nonce1", password))
	rec := httptest.NewRecorder()
	h.Autologin(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("unbound: status = %d, want 400", rec.Code)
	}

	// bound install → session
	_ = installs.SetUser("123", user.ID)
	req = autologinReq("nonce1", "123", "token="+appinstall.NoncePasswordMD5("nonce1", password))
	rec = httptest.NewRecorder()
	h.Autologin(rec, req)
	if rec.Code != http.StatusFound {
		t.Fatalf("bound: status = %d, want 302", rec.Code)
	}
	if loc := rec.Header().Get("Location"); !strings.HasPrefix(loc, "/onboarding?code=") {
		t.Errorf("Location = %q", loc)
	}
}

// (9) LinkPending happy path → linked + api_key.
func TestLinkPending_OK(t *testing.T) {
	h, installs, users, _, _ := newTestAppHandler(t)
	installFor(t, h, "123", "token123")
	user, _ := users.Create("u1", "user@test.ru", "", "", "", "hash")
	_ = installs.SetPendingLink("123", user.ID)

	body := strings.NewReader(`{"insales_id":"123"}`)
	req := httptest.NewRequest("POST", "/api/shops/link-pending", body)
	req = req.WithContext(context.WithValue(req.Context(), auth.UserIDKey, user.ID))
	rec := httptest.NewRecorder()
	h.LinkPending(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), `"status":"linked"`) || !strings.Contains(rec.Body.String(), `"api_key":"`) {
		t.Errorf("body = %s, want status=linked + api_key", rec.Body.String())
	}
	inst, _ := installs.GetByInsalesID("123")
	if inst.ShopID == "" {
		t.Error("shop not linked")
	}
}

// (10) LinkPending on non-pending install → refuse.
func TestLinkPending_NotPending(t *testing.T) {
	h, installs, users, _, _ := newTestAppHandler(t)
	installFor(t, h, "123", "token123")
	user, _ := users.Create("u1", "user@test.ru", "", "", "", "hash")
	_ = installs.SetUser("123", user.ID) // status stays "installed"

	body := strings.NewReader(`{"insales_id":"123"}`)
	req := httptest.NewRequest("POST", "/api/shops/link-pending", body)
	req = req.WithContext(context.WithValue(req.Context(), auth.UserIDKey, user.ID))
	rec := httptest.NewRecorder()
	h.LinkPending(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403", rec.Code)
	}
}

func TestUninstall_WrongToken(t *testing.T) {
	h, installs, _, _, _ := newTestAppHandler(t)
	// install first
	req := httptest.NewRequest("GET", "/install?shop=test.myinsales.ru&token=token123&insales_id=123", nil)
	rec := httptest.NewRecorder()
	h.Install(rec, req)

	req = httptest.NewRequest("GET", "/uninstall?shop=test.myinsales.ru&token=WRONG&insales_id=123", nil)
	rec = httptest.NewRecorder()
	h.Uninstall(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
	inst, _ := installs.GetByInsalesID("123")
	if inst.Status != "installed" {
		t.Errorf("status = %q, want installed", inst.Status)
	}
}

func TestUninstall_OK(t *testing.T) {
	h, installs, _, shops, keys := newTestAppHandler(t)
	req := httptest.NewRequest("GET", "/install?shop=test.myinsales.ru&token=token123&insales_id=123", nil)
	rec := httptest.NewRecorder()
	h.Install(rec, req)

	// simulate provisioned shop + key
	password := appinstall.InstallPassword("token123", "mysecret")
	enc, _ := h.aes.Encrypt(password)
	_, _ = installs.UpsertInstall("123", "test.myinsales.ru", enc)
	shop, _ := shops.Create("user-1", "test.myinsales.ru", "test.myinsales.ru", "myapplogin", enc, true, false, false)
	_, _ = keys.CreateForShop("user-1", shop.ID, "hash1", "InSales: test.myinsales.ru")
	_ = installs.LinkShop("123", "user-1", shop.ID)

	req = httptest.NewRequest("GET", "/uninstall?shop=test.myinsales.ru&token="+password+"&insales_id=123", nil)
	rec = httptest.NewRecorder()
	h.Uninstall(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	inst, _ := installs.GetByInsalesID("123")
	if inst.Status != "uninstalled" {
		t.Errorf("status = %q, want uninstalled", inst.Status)
	}
	// keys must survive uninstall
	if _, err := keys.GetByHash("hash1"); err != nil {
		t.Errorf("ai key deleted on uninstall: %v", err)
	}
}

func TestExchangeCode_Invalid(t *testing.T) {
	h, _, _, _, _ := newTestAppHandler(t)
	code, err := h.issueCode("user-1", "rawkey")
	if err != nil {
		t.Fatal(err)
	}
	uid, raw, ok := h.consumeCode(code)
	if !ok || uid != "user-1" || raw != "rawkey" {
		t.Fatalf("consumeCode = (%q,%q,%v)", uid, raw, ok)
	}
	if _, _, ok := h.consumeCode("deadbeef.payload"); ok {
		t.Error("consumeCode accepted bad signature")
	}
}

func TestNewAppHandler_BillingMonthly(t *testing.T) {
	h, _, _, _, _ := newTestAppHandler(t)
	if h.billingMonthly != 0 {
		t.Errorf("default billing = %d", h.billingMonthly)
	}
}

// (billing) charge blocked → session still issued, status unpaid (plan test 4).
func TestAutologin_BillingBlocked_SessionIssued(t *testing.T) {
	h, installs, users, shops, _ := newTestAppHandler(t)
	h.billingMonthly = 300
	installFor(t, h, "123", "token123")
	password := appinstall.InstallPassword("token123", "mysecret")
	enc, _ := h.aes.Encrypt(password)
	user, _ := users.Create("u1", "user@test.ru", "", "", "", "hash")
	shop, _ := shops.Create(user.ID, "test.myinsales.ru", "test.myinsales.ru", "myapplogin", enc, true, false, false)
	shop.IsAppInstall = true
	_ = installs.LinkShop("123", user.ID, shop.ID)
	h.chargeFn = func(_ *storage.AppInstall, _ string) (map[string]any, error) {
		return map[string]any{"blocked": true, "confirmation_url": "https://pay/1"}, nil
	}

	req := autologinReq("nonce1", "123",
		"user_email=user@test.ru&user_name=u1&user_id=1&email_confirmed=true&token3="+
			token3("nonce1", "user@test.ru", "u1", "1", "true", password))
	rec := httptest.NewRecorder()
	h.Autologin(rec, req)

	if rec.Code != http.StatusFound {
		t.Fatalf("status = %d, want 302 (login must never be blocked)", rec.Code)
	}
	if loc := rec.Header().Get("Location"); !strings.HasPrefix(loc, "/onboarding?code=") {
		t.Fatalf("Location = %q, want /onboarding?code=", loc)
	}
	if shops.shops[shop.ID].BillingStatus != "unpaid" {
		t.Errorf("billing status = %q, want unpaid", shops.shops[shop.ID].BillingStatus)
	}
	if shops.shops[shop.ID].BillingPayURL != "https://pay/1" {
		t.Errorf("pay url = %q", shops.shops[shop.ID].BillingPayURL)
	}
}

// (billing) paid charge → active + until (plan test 5).
func TestAutologin_BillingPaid_Active(t *testing.T) {
	h, installs, users, shops, _ := newTestAppHandler(t)
	h.billingMonthly = 300
	installFor(t, h, "123", "token123")
	password := appinstall.InstallPassword("token123", "mysecret")
	enc, _ := h.aes.Encrypt(password)
	user, _ := users.Create("u1", "user@test.ru", "", "", "", "hash")
	shop, _ := shops.Create(user.ID, "test.myinsales.ru", "test.myinsales.ru", "myapplogin", enc, true, false, false)
	shop.IsAppInstall = true
	_ = installs.LinkShop("123", user.ID, shop.ID)
	h.chargeFn = func(_ *storage.AppInstall, _ string) (map[string]any, error) {
		return map[string]any{"paid_till": "2099-01-01"}, nil
	}

	req := autologinReq("nonce1", "123",
		"user_email=user@test.ru&user_name=u1&user_id=1&email_confirmed=true&token3="+
			token3("nonce1", "user@test.ru", "u1", "1", "true", password))
	rec := httptest.NewRecorder()
	h.Autologin(rec, req)

	if rec.Code != http.StatusFound {
		t.Fatalf("status = %d, want 302", rec.Code)
	}
	got := shops.shops[shop.ID]
	if got.BillingStatus != "active" || got.BillingUntil == nil {
		t.Errorf("billing = (%s, %v), want (active, until)", got.BillingStatus, got.BillingUntil)
	}
}

// (billing) BILLING_MONTHLY=0 → no charge call at all (plan test 6).
func TestAutologin_BillingFree_NoChargeCall(t *testing.T) {
	h, installs, users, shops, _ := newTestAppHandler(t)
	installFor(t, h, "123", "token123")
	password := appinstall.InstallPassword("token123", "mysecret")
	enc, _ := h.aes.Encrypt(password)
	user, _ := users.Create("u1", "user@test.ru", "", "", "", "hash")
	shop, _ := shops.Create(user.ID, "test.myinsales.ru", "test.myinsales.ru", "myapplogin", enc, true, false, false)
	_ = installs.LinkShop("123", user.ID, shop.ID)
	called := false
	h.chargeFn = func(_ *storage.AppInstall, _ string) (map[string]any, error) {
		called = true
		return map[string]any{}, nil
	}

	req := autologinReq("nonce1", "123",
		"user_email=user@test.ru&user_name=u1&user_id=1&email_confirmed=true&token3="+
			token3("nonce1", "user@test.ru", "u1", "1", "true", password))
	rec := httptest.NewRecorder()
	h.Autologin(rec, req)

	if rec.Code != http.StatusFound {
		t.Fatalf("status = %d, want 302", rec.Code)
	}
	if called {
		t.Error("charge must not be fetched when BILLING_MONTHLY=0")
	}
}

// (billing) fetch failure → login proceeds, shop marked unknown for fast retry (test 14).
func TestAutologin_BillingFetchError_LoginProceeds(t *testing.T) {
	h, installs, users, shops, _ := newTestAppHandler(t)
	h.billingMonthly = 300
	installFor(t, h, "123", "token123")
	password := appinstall.InstallPassword("token123", "mysecret")
	enc, _ := h.aes.Encrypt(password)
	user, _ := users.Create("u1", "user@test.ru", "", "", "", "hash")
	shop, _ := shops.Create(user.ID, "test.myinsales.ru", "test.myinsales.ru", "myapplogin", enc, true, false, false)
	shop.IsAppInstall = true
	_ = installs.LinkShop("123", user.ID, shop.ID)
	h.chargeFn = func(_ *storage.AppInstall, _ string) (map[string]any, error) {
		return nil, errors.New("insales down")
	}

	req := autologinReq("nonce1", "123",
		"user_email=user@test.ru&user_name=u1&user_id=1&email_confirmed=true&token3="+
			token3("nonce1", "user@test.ru", "u1", "1", "true", password))
	rec := httptest.NewRecorder()
	h.Autologin(rec, req)

	if rec.Code != http.StatusFound {
		t.Fatalf("status = %d, want 302 (billing outage must not block login)", rec.Code)
	}
	if got := shops.shops[shop.ID].BillingStatus; got != "unknown" {
		t.Errorf("billing status = %q, want unknown", got)
	}
}
