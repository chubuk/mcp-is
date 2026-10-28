package handlers

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/chubuk/mcp-is/internal/appinstall"
	"github.com/chubuk/mcp-is/internal/auth"
	"github.com/chubuk/mcp-is/internal/core"
	"github.com/chubuk/mcp-is/internal/integration/insales"
	"github.com/chubuk/mcp-is/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

// AppHandler implements the InSales external-application callbacks:
// install, uninstall, autologin and one-time code exchange for the SPA.
type AppHandler struct {
	installs storage.AppInstallRepositoryInterface
	users    storage.UserRepositoryInterface
	shops    storage.ShopRepositoryInterface
	aiKeys   storage.AIKeyRepositoryInterface
	aes      *auth.AESCipher
	jwt      *auth.JWTManager
	manager  *core.AccountManager

	appKey         string
	secretKey      string
	domain         string
	hmacSecret     string
	billingMonthly int // 0 = free (default)

	// chargeFn overrides the InSales charge fetch in tests (nil = real API).
	chargeFn func(inst *storage.AppInstall, password string) (map[string]any, error)
}

func NewAppHandler(
	installs storage.AppInstallRepositoryInterface,
	users storage.UserRepositoryInterface,
	shops storage.ShopRepositoryInterface,
	aiKeys storage.AIKeyRepositoryInterface,
	aes *auth.AESCipher,
	jwt *auth.JWTManager,
	manager *core.AccountManager,
	appKey, secretKey, domain, hmacSecret string,
	billingMonthly int,
) *AppHandler {
	return &AppHandler{
		installs: installs, users: users, shops: shops, aiKeys: aiKeys,
		aes: aes, jwt: jwt, manager: manager,
		appKey: appKey, secretKey: secretKey, domain: domain, hmacSecret: hmacSecret,
		billingMonthly: billingMonthly,
	}
}

func (h *AppHandler) enabled() bool {
	return h.appKey != "" && h.secretKey != ""
}

func appBadRequest(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(http.StatusBadRequest)
	_, _ = w.Write([]byte("bad request"))
}

// Install handles GET /install?shop=&token=&insales_id=
// password = MD5(token + secret_key); respond 200 only after the password is stored.
func (h *AppHandler) Install(w http.ResponseWriter, r *http.Request) {
	if !h.enabled() {
		http.NotFound(w, r)
		return
	}
	q := r.URL.Query()
	shop, token, insalesID := q.Get("shop"), q.Get("token"), q.Get("insales_id")
	if shop == "" || token == "" || insalesID == "" {
		slog.Warn("app install: missing params", "shop", shop, "insales_id", insalesID)
		appBadRequest(w)
		return
	}
	shop = normalizeShopHost(shop)
	if shop == "" {
		appBadRequest(w)
		return
	}

	password := appinstall.InstallPassword(token, h.secretKey)
	enc, err := h.aes.Encrypt(password)
	if err != nil {
		slog.Error("app install: encrypt", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	prev, _ := h.installs.GetByInsalesID(insalesID)
	if prev != nil && prev.Status == "installed" {
		slog.Warn("app install: overwrite existing install", "insales_id", insalesID, "shop", shop)
	}

	inst, err := h.installs.UpsertInstall(insalesID, shop, enc)
	if err != nil {
		slog.Error("app install: upsert", "error", err, "insales_id", insalesID)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// reinstall rotates the password — refresh the shops copy the API client reads,
	// then drop the cached client. Order matters: write first, invalidate last,
	// non-200 on failure so InSales retries instead of diverging rows.
	if inst.ShopID != "" {
		if err := h.shops.UpdateAPIPassword(inst.ShopID, enc); err != nil {
			slog.Error("app install: refresh shop password", "error", err, "insales_id", insalesID, "shop_id", inst.ShopID)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		h.manager.InvalidateShop(inst.ShopID)
	}

	slog.Info("app installed", "insales_id", insalesID, "shop", shop)
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}

// Uninstall handles GET /uninstall?shop=&token=<password>&insales_id=
func (h *AppHandler) Uninstall(w http.ResponseWriter, r *http.Request) {
	if !h.enabled() {
		http.NotFound(w, r)
		return
	}
	q := r.URL.Query()
	shop, token, insalesID := q.Get("shop"), q.Get("token"), q.Get("insales_id")
	if shop == "" || token == "" || insalesID == "" {
		appBadRequest(w)
		return
	}

	inst, err := h.installs.GetByInsalesID(insalesID)
	if err != nil {
		slog.Warn("app uninstall: unknown install", "insales_id", insalesID)
		appBadRequest(w)
		return
	}
	password, err := h.aes.Decrypt(inst.APIPasswordEncrypted)
	if err != nil {
		slog.Error("app uninstall: decrypt", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if subtle.ConstantTimeCompare([]byte(password), []byte(token)) != 1 {
		slog.Warn("app uninstall: bad token", "insales_id", insalesID)
		appBadRequest(w)
		return
	}

	if inst.ShopID != "" {
		if err := h.shops.Delete(inst.ShopID); err != nil {
			slog.Warn("app uninstall: delete shop", "error", err, "shop_id", inst.ShopID)
		}
		h.manager.InvalidateShop(inst.ShopID)
	}
	if err := h.installs.MarkUninstalled(insalesID); err != nil {
		slog.Error("app uninstall: mark", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	slog.Info("app uninstalled", "insales_id", insalesID, "shop", normalizeShopHost(shop))
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}

// AppLanding handles GET /app — entry page from the InSales admin «URL входа».
func (h *AppHandler) AppLanding(w http.ResponseWriter, r *http.Request) {
	if !h.enabled() {
		http.NotFound(w, r)
		return
	}
	q := r.URL.Query()
	shop, insalesID := normalizeShopHost(q.Get("shop")), q.Get("insales_id")
	if shop == "" || insalesID == "" {
		appBadRequest(w)
		return
	}

	// Always do the InSales round-trip: the SPA JWT lives in localStorage and is
	// invisible to the server, so a shortcut would drop users on a clean browser.
	nonce, err := randomHex(16)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Cookie binds nonce to the shop — /autologin does not receive shop/insales_id in query.
	http.SetCookie(w, &http.Cookie{
		Name:     "autologin_nonce",
		Value:    nonce + ":" + insalesID,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   300,
	})

	loginURL := fmt.Sprintf("https://%s/autologin", h.domain)
	target := fmt.Sprintf("http://%s/admin/applications/%s/login?token=%s&login=%s",
		shop, h.appKey, nonce, loginURL)
	http.Redirect(w, r, target, http.StatusFound)
}

// Autologin handles GET /autologin?token3=&user_email=&user_name=&user_id=&email_сonfirmed=
func (h *AppHandler) Autologin(w http.ResponseWriter, r *http.Request) {
	if !h.enabled() {
		http.NotFound(w, r)
		return
	}
	cookie, err := r.Cookie("autologin_nonce")
	if err != nil {
		slog.Warn("app autologin: no nonce cookie")
		appBadRequest(w)
		return
	}
	nonce, insalesID, ok := strings.Cut(cookie.Value, ":")
	if !ok || nonce == "" || insalesID == "" {
		appBadRequest(w)
		return
	}

	inst, err := h.installs.GetByInsalesID(insalesID)
	if err != nil {
		slog.Warn("app autologin: unknown install", "insales_id", insalesID)
		appBadRequest(w)
		return
	}
	password, err := h.aes.Decrypt(inst.APIPasswordEncrypted)
	if err != nil {
		slog.Error("app autologin: decrypt", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	q := r.URL.Query()
	// InSales docs use a Cyrillic «с» in email_сonfirmed — accept both spellings.
	emailConfirmed := q.Get("email_confirmed")
	if emailConfirmed == "" {
		emailConfirmed = q.Get("email_сonfirmed")
	}
	userEmail, userName, userID := q.Get("user_email"), q.Get("user_name"), q.Get("user_id")
	token3 := q.Get("token3")

	if token3 != "" {
		if !appinstall.ValidateToken3(nonce, userEmail, userName, userID, emailConfirmed, password, token3) {
			slog.Warn("app autologin: bad token3", "insales_id", insalesID)
			appBadRequest(w)
			return
		}
		if emailConfirmed != "true" {
			slog.Warn("app autologin: email not confirmed", "insales_id", insalesID)
			http.Error(w, "email not confirmed", http.StatusForbidden)
			return
		}
	} else if token := q.Get("token"); token != "" {
		// Short vendor format (?token=MD5(nonce+password)) carries no user identity —
		// usable only for an install already bound to a user.
		if inst.UserID == "" || subtle.ConstantTimeCompare([]byte(strings.ToLower(token)),
			[]byte(appinstall.NoncePasswordMD5(nonce, password))) != 1 {
			slog.Warn("app autologin: short token without user binding", "insales_id", insalesID)
			appBadRequest(w)
			return
		}
		userEmail, userName = "", ""
	} else {
		slog.Warn("app autologin: no token3/token", "insales_id", insalesID)
		appBadRequest(w)
		return
	}

	// Expire the nonce — single use (cookie-only; replay window 300s accepted by design).
	http.SetCookie(w, &http.Cookie{Name: "autologin_nonce", Value: "", Path: "/", MaxAge: -1, HttpOnly: true, Secure: true, SameSite: http.SameSiteLaxMode})

	existing, err := h.users.GetByEmail(userEmail)
	if err == nil && existing != nil {
		switch {
		case inst.UserID == existing.ID:
			// Install already bound to this user: no password prompt, no duplicate
			// shop/key — provisionUserAndShop is idempotent (recreates shop after
			// uninstall, reuses it otherwise; rawKey is "" when nothing is created).
			user, rawKey, perr := h.provisionUserAndShop(inst, userEmail, userName)
			if perr != nil {
				slog.Error("app autologin: re-provision", "error", perr, "insales_id", insalesID)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			if h.billingMonthly > 0 {
				if berr := h.syncBilling(inst, password); berr != nil {
					slog.Warn("app autologin: billing check skipped", "error", berr, "insales_id", insalesID)
				}
			}
			code, cerr := h.issueCode(user.ID, rawKey)
			if cerr != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			slog.Info("app autologin: returning user", "insales_id", insalesID, "user_id", user.ID)
			http.Redirect(w, r, "/onboarding?code="+code, http.StatusFound)
			return
		case inst.UserID != "" && inst.UserID != existing.ID:
			// Email reassigned to another account — never rewrite ownership.
			slog.Warn("app autologin: user conflict, refusing link", "insales_id", insalesID,
				"inst_user_id", inst.UserID, "email_user_id", existing.ID)
			http.Error(w, "account conflict", http.StatusConflict)
			return
		default:
			// Anti-takeover: never issue a session for a pre-existing account.
			if err := h.installs.SetPendingLink(insalesID, existing.ID); err != nil {
				slog.Error("app autologin: pending link", "error", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			slog.Info("app autologin: pending link to existing user", "insales_id", insalesID, "user_id", existing.ID)
			http.Redirect(w, r, "/login?pending_install="+insalesID, http.StatusFound)
			return
		}
	}

	user, rawKey, err := h.provisionUserAndShop(inst, userEmail, userName)
	if err != nil {
		slog.Error("app autologin: provision", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Billing never blocks login: the status feeds the UI banner and the MCP
	// gate. Charge check failure only delays the check to the lazy TTL retry.
	if h.billingMonthly > 0 {
		if berr := h.syncBilling(inst, password); berr != nil {
			slog.Warn("app autologin: billing check skipped", "error", berr, "insales_id", insalesID)
		}
	}

	code, err := h.issueCode(user.ID, rawKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/onboarding?code="+code, http.StatusFound)
}

// ExchangeCode handles POST /api/exchange-code — swaps the one-time code for a JWT (+ fresh MCP key).
func (h *AppHandler) ExchangeCode(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Code string `json:"code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Code == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request"})
		return
	}
	userID, rawKey, ok := h.consumeCode(req.Code)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid or expired code"})
		return
	}
	user, err := h.users.GetByID(userID)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid code"})
		return
	}
	token, err := h.jwt.GenerateToken(user.ID, user.Nickname, user.Role)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{
		"token": token, "user_id": user.ID, "nickname": user.Nickname, "role": user.Role,
		"api_key": rawKey,
	})
}

// LinkPending handles POST /api/shops/link-pending — completes hybrid link after login.
func (h *AppHandler) LinkPending(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(auth.UserIDKey).(string)
	var req struct {
		InsalesID string `json:"insales_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.InsalesID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request"})
		return
	}
	inst, err := h.installs.GetByInsalesID(req.InsalesID)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "установка не найдена"})
		return
	}
	if inst.Status != "pending_link" || inst.UserID != userID {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "привязка недоступна"})
		return
	}
	_, rawKey, err := h.provisionUserAndShop(inst, "", "")
	if err != nil {
		slog.Error("link pending: provision", "error", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"status": "linked", "api_key": rawKey})
}

// --- provisioning helpers ---

// provisionUserAndShop creates (or reuses) the user, upserts the shop and a shop-bound MCP key.
// email/userName are only used when creating a new user (empty on link-pending path).

// ensureRecurringCharge creates or re-checks the singleton RecurringApplicationCharge.
// Returns the raw charge map for core.MapCharge.
func (h *AppHandler) ensureRecurringCharge(inst *storage.AppInstall, password string) (map[string]any, error) {
	client, err := insales.NewClient(inst.ShopHost, h.appKey, password)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()

	existing, err := client.GetRecurringApplicationCharge(ctx, "")
	if err == nil && existing != nil {
		if m, ok := existing.(map[string]any); ok {
			return m, nil
		}
	}

	charge, err := client.CreateRecurringApplicationCharge(ctx, map[string]any{
		"monthly": h.billingMonthly,
	})
	if err != nil {
		return nil, fmt.Errorf("create recurring charge: %w", err)
	}
	m, _ := charge.(map[string]any)
	return m, nil
}

// syncBilling refreshes the charge state at login — at most once per TTL window.
// Never blocks login: on failure the shop is marked unknown for a fast lazy retry.
func (h *AppHandler) syncBilling(inst *storage.AppInstall, password string) error {
	if h.billingMonthly <= 0 || inst.ShopID == "" {
		return nil
	}
	shop, err := h.shops.GetByID(inst.ShopID)
	if err != nil {
		return err
	}
	if !core.NeedsBillingRefresh(shop, time.Now()) {
		return nil
	}
	fetch := h.ensureRecurringCharge
	if h.chargeFn != nil {
		fetch = h.chargeFn
	}
	charge, err := fetch(inst, password)
	if err != nil {
		// keep the gate fail-open and retry lazily within the hour
		if uerr := h.shops.UpdateBillingStatus(inst.ShopID, core.BillingUnknown, nil, ""); uerr != nil {
			return uerr
		}
		return err
	}
	status, until, payURL := core.MapCharge(charge, time.Now())
	return h.shops.UpdateBillingStatus(inst.ShopID, status, until, payURL)
}

func (h *AppHandler) provisionUserAndShop(inst *storage.AppInstall, email, userName string) (*storage.User, string, error) {
	userID := inst.UserID
	if userID == "" {
		nickname := uniqueNickname(userName, email)
		hash, err := bcrypt.GenerateFromPassword([]byte(randomPassword()), bcrypt.DefaultCost)
		if err != nil {
			return nil, "", err
		}
		user, err := h.users.Create(nickname, email, "", "", "", string(hash))
		if err != nil {
			return nil, "", fmt.Errorf("create user: %w", err)
		}
		if err := h.users.UpdateRequirePasswordChange(user.ID, true); err != nil {
			slog.Warn("require_password_change", "error", err)
		}
		userID = user.ID
	}

	shop, rawKey := inst.ShopID, ""
	if shop != "" {
		// Returning user / already provisioned — reuse shop, don't mint another key.
		if _, err := h.shops.GetByID(shop); err != nil {
			return nil, "", fmt.Errorf("get shop: %w", err)
		}
	} else {
		created, err := h.shops.Create(userID, inst.ShopHost, inst.ShopHost, h.appKey, inst.APIPasswordEncrypted, true, false, false)
		if err != nil {
			return nil, "", fmt.Errorf("create shop: %w", err)
		}
		shop = created.ID
		key, keyHash, err := auth.GenerateAPIKey()
		if err != nil {
			return nil, "", err
		}
		if _, err := h.aiKeys.CreateForShop(userID, shop, keyHash, "InSales: "+inst.ShopHost); err != nil {
			return nil, "", fmt.Errorf("create ai key: %w", err)
		}
		if err := h.installs.LinkShop(inst.InsalesID, userID, shop); err != nil {
			return nil, "", err
		}
		rawKey = key
	}
	user, err := h.users.GetByID(userID)
	if err != nil {
		return nil, "", err
	}
	inst.ShopID = shop // LinkShop writes DB only — keep the in-memory install in sync
	return user, rawKey, nil
}

// --- one-time codes (HMAC, short TTL). Payload: userID:exp:rawKey ---

func (h *AppHandler) issueCode(userID, rawKey string) (string, error) {
	exp := time.Now().Add(2 * time.Minute).Unix()
	payload := fmt.Sprintf("%s:%d:%s", userID, exp, rawKey)
	mac := hmac.New(sha256.New, []byte(h.hmacSecret))
	mac.Write([]byte(payload))
	return hex.EncodeToString(mac.Sum(nil)) + "." + payload, nil
}

func (h *AppHandler) consumeCode(code string) (string, string, bool) {
	sig, payload, ok := strings.Cut(code, ".")
	if !ok {
		return "", "", false
	}
	mac := hmac.New(sha256.New, []byte(h.hmacSecret))
	mac.Write([]byte(payload))
	want := hex.EncodeToString(mac.Sum(nil))
	if subtle.ConstantTimeCompare([]byte(sig), []byte(want)) != 1 {
		return "", "", false
	}
	parts := strings.SplitN(payload, ":", 3)
	if len(parts) != 3 {
		return "", "", false
	}
	var exp int64
	if _, err := fmt.Sscanf(parts[1], "%d", &exp); err != nil || time.Now().Unix() > exp {
		return "", "", false
	}
	return parts[0], parts[2], true
}

func normalizeShopHost(shop string) string {
	shop = strings.TrimSpace(strings.ToLower(shop))
	shop = strings.TrimPrefix(shop, "https://")
	shop = strings.TrimPrefix(shop, "http://")
	shop = strings.TrimSuffix(shop, "/")
	if shop == "" || strings.ContainsAny(shop, " ?#&\t\n") {
		return ""
	}
	return shop
}

func randomHex(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func randomPassword() string {
	s, _ := randomHex(16)
	return "!" + s + "Aa1"
}

func uniqueNickname(userName, email string) string {
	base := strings.TrimSpace(userName)
	if base == "" {
		base = strings.SplitN(email, "@", 2)[0]
	}
	base = strings.Map(func(r rune) rune {
		if r < 32 {
			return -1
		}
		return r
	}, base)
	if base == "" {
		base = "user"
	}
	// uniqueness is enforced by DB; suffix reduces collisions on create
	suffix, _ := randomHex(3)
	return base + "-" + suffix[:5]
}
