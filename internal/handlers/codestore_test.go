package handlers

import "testing"

func TestCodeStore_IssueVerify(t *testing.T) {
	s := newCodeStore()
	code, err := s.issue("register:u1")
	if err != nil {
		t.Fatalf("issue: %v", err)
	}
	if len(code) != 6 {
		t.Errorf("want 6-digit code, got %q", code)
	}
	if err := s.verify("register:u1", code); err != nil {
		t.Fatalf("verify: %v", err)
	}
	// single-use: повторный verify того же кода — отказ
	if err := s.verify("register:u1", code); err == nil {
		t.Error("replay must fail")
	}
}

func TestCodeStore_WrongAttemptsBurn(t *testing.T) {
	s := newCodeStore()
	code, _ := s.issue("link:1")
	for i := 0; i < codeMaxTries; i++ {
		if err := s.verify("link:1", "000000"); err == nil {
			t.Fatal("wrong code must fail")
		}
	}
	// после codeMaxTries неверных попыток код сожжён — даже верный не пройдёт
	if err := s.verify("link:1", code); err == nil {
		t.Error("burned code must not verify")
	}
}

func TestCodeStore_Cooldown(t *testing.T) {
	s := newCodeStore()
	if _, err := s.issue("register:u2"); err != nil {
		t.Fatalf("first issue: %v", err)
	}
	if _, err := s.issue("register:u2"); err != ErrCodeCooldown {
		t.Errorf("want ErrCodeCooldown, got %v", err)
	}
	// другой ключ не задет
	if _, err := s.issue("register:u3"); err != nil {
		t.Errorf("other key must issue, got %v", err)
	}
}

func TestCodeStore_UnknownKey(t *testing.T) {
	s := newCodeStore()
	if err := s.verify("register:none", "123456"); err != ErrCodeInvalid {
		t.Errorf("want ErrCodeInvalid, got %v", err)
	}
}
