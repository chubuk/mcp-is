package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"math/big"
	"sync"
	"time"
)

var (
	ErrCodeCooldown = errors.New("code resend cooldown")
	ErrCodeInvalid  = errors.New("invalid or expired code")
)

const (
	codeTTL      = 10 * time.Minute
	codeCooldown = 90 * time.Second
	codeMaxTries = 5
)

type codeEntry struct {
	code     string
	expires  time.Time
	lastSent time.Time
	tries    int
}

// codeStore holds single-use email OTPs (purposes: "register", "link") in memory.
// Single-VPS process: codes lost on restart are simply re-requested.
type codeStore struct {
	mu sync.Mutex
	m  map[string]*codeEntry
}

func newCodeStore() *codeStore {
	return &codeStore{m: make(map[string]*codeEntry)}
}

// issue generates a 6-digit code for key, enforcing a resend cooldown.
func (s *codeStore) issue(key string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	if e, ok := s.m[key]; ok && now.Sub(e.lastSent) < codeCooldown {
		return "", ErrCodeCooldown
	}
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", err
	}
	code := pad6(n.Int64())
	s.m[key] = &codeEntry{code: code, expires: now.Add(codeTTL), lastSent: now}
	return code, nil
}

// verify consumes the code for key: single-use, TTL-limited, max tries.
func (s *codeStore) verify(key, code string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.m[key]
	if !ok || time.Now().After(e.expires) {
		delete(s.m, key)
		return ErrCodeInvalid
	}
	e.tries++
	if e.tries > codeMaxTries {
		delete(s.m, key)
		return ErrCodeInvalid
	}
	if e.code != code {
		if e.tries >= codeMaxTries {
			delete(s.m, key)
		}
		return ErrCodeInvalid
	}
	delete(s.m, key)
	return nil
}

func pad6(n int64) string {
	s := [6]byte{'0', '0', '0', '0', '0', '0'}
	i := 5
	for n > 0 {
		s[i] = byte('0' + n%10)
		n /= 10
		i--
	}
	return string(s[:])
}

// --- onboarding one-time exchange codes (SPA autologin handoff) ---

const onboardTTL = 2 * time.Minute

type onboardEntry struct {
	userID  string
	rawKey  string
	expires time.Time
}

// onboardStore holds single-use onboarding payloads keyed by the code itself
// (inverted vs codeStore: code → payload, not key → code). Single-VPS process:
// codes lost on restart are simply re-requested via a fresh autologin.
type onboardStore struct {
	mu sync.Mutex
	m  map[string]*onboardEntry
}

func newOnboardStore() *onboardStore {
	return &onboardStore{m: make(map[string]*onboardEntry)}
}

// issue stores the payload under a fresh random code (TTL 2 min).
func (s *onboardStore) issue(userID, rawKey string) (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	code := hex.EncodeToString(b)
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m[code] = &onboardEntry{userID: userID, rawKey: rawKey, expires: time.Now().Add(onboardTTL)}
	return code, nil
}

// take atomically consumes the code: concurrent takes yield exactly one winner.
func (s *onboardStore) take(code string) (string, string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.m[code]
	if !ok || time.Now().After(e.expires) {
		delete(s.m, code)
		return "", "", false
	}
	delete(s.m, code)
	return e.userID, e.rawKey, true
}
