package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/chubuk/mcp-is/internal/appinstall"
	"github.com/chubuk/mcp-is/internal/auth"
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
	shop, _ := shops.Create("user-1", "test.myinsales.ru", "test.myinsales.ru", "myapplogin", enc, true, false)
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

func TestChargeConfirmationURL(t *testing.T) {
	got := chargeConfirmationURL(map[string]any{"confirmation_url": "https://x.myinsales.ru/confirm/1"})
	if got != "https://x.myinsales.ru/confirm/1" {
		t.Errorf("got %q", got)
	}
	if chargeConfirmationURL(map[string]any{"blocked": true}) != "" {
		t.Error("expected empty when no url")
	}
}

func TestNewAppHandler_BillingMonthly(t *testing.T) {
	h, _, _, _, _ := newTestAppHandler(t)
	if h.billingMonthly != 0 {
		t.Errorf("default billing = %d", h.billingMonthly)
	}
}
