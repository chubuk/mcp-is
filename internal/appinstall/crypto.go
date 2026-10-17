// Package appinstall implements the InSales external-application install protocol
// (password = MD5(token + secret_key), autologin token3).
package appinstall

import (
	"crypto/md5"
	"crypto/subtle"
	"encoding/hex"
	"strings"
)

// InstallPassword derives the per-shop app password from the one-time install token.
// password = MD5(token + secret_key), hex lowercase (InSales protocol).
func InstallPassword(token, secretKey string) string {
	sum := md5.Sum([]byte(token + secretKey))
	return hex.EncodeToString(sum[:])
}

// ValidateToken3 checks the autologin redirect token.
// token3 = MD5(nonce + user_email + user_name + user_id + email_confirmed + password)
// where password is the per-shop install password (not secret_key).
// emailConfirmed is the raw query value ("true"/"false").
func ValidateToken3(nonce, userEmail, userName, userID, emailConfirmed, password, token3 string) bool {
	sum := md5.Sum([]byte(nonce + userEmail + userName + userID + emailConfirmed + password))
	want := hex.EncodeToString(sum[:])
	return subtle.ConstantTimeCompare([]byte(strings.ToLower(token3)), []byte(want)) == 1
}
