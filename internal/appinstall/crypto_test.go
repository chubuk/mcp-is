package appinstall

import "testing"

// Vector from docs/api/api_and_apps.md: MD5("token123mysecret") = 4c3f4c197336eb97475506e71c839c71
func TestInstallPassword_DocVector(t *testing.T) {
	got := InstallPassword("token123", "mysecret")
	want := "4c3f4c197336eb97475506e71c839c71"
	if got != want {
		t.Errorf("InstallPassword() = %q, want %q", got, want)
	}
}

func TestValidateToken3_OK(t *testing.T) {
	// password from the doc vector
	password := InstallPassword("token123", "mysecret")
	// token3 = MD5(nonce + user_email + user_name + user_id + email_confirmed + password)
	// computed via the same helper path: recompute expected with a known nonce
	nonce := "3ad376b668c48c15d4beed467a567ef7"
	// Build expected by calling InstallPassword on the concatenated string
	// (MD5 of concatenation — same as token3 definition).
	token3 := InstallPassword(nonce+"user@myshop.ru"+"igor"+"1"+"true", password)

	if !ValidateToken3(nonce, "user@myshop.ru", "igor", "1", "true", password, token3) {
		t.Error("ValidateToken3: expected true for valid token3")
	}
}

func TestValidateToken3_BadToken(t *testing.T) {
	password := InstallPassword("token123", "mysecret")
	if ValidateToken3("nonce", "user@myshop.ru", "igor", "1", "true", password, "deadbeef") {
		t.Error("ValidateToken3: expected false for wrong token3")
	}
}

func TestValidateToken3_BadPassword(t *testing.T) {
	token3 := InstallPassword("nonce"+"user@myshop.ru"+"igor"+"1"+"true", "right-password")
	if ValidateToken3("nonce", "user@myshop.ru", "igor", "1", "true", "wrong-password", token3) {
		t.Error("ValidateToken3: expected false when password differs")
	}
}
