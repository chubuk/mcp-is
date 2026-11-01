package pdn

import (
	"strings"
	"testing"
)

func TestMaskCSV(t *testing.T) {
	in := "id,number,total_price,client_name,created_at\n" +
		"1,1000,1920.00,\"Иванов Иван\",2026-11-01\n"
	out := MaskCSV(in)
	if strings.Contains(out, "Иванов") {
		t.Errorf("client_name leaked: %q", out)
	}
	if !strings.Contains(out, "И***") {
		t.Errorf("mask missing: %q", out)
	}
	if !strings.Contains(out, "1920.00") {
		t.Errorf("non-PII column modified: %q", out)
	}
}

func TestMaskCSV_ClientsColumns(t *testing.T) {
	in := "id,name,surname,email,phone,created_at\n" +
		"1,Иван,Петров,ivan@example.com,+79991234567,2026-11-01\n"
	out := MaskCSV(in)
	for _, leak := range []string{"Иван", "Петров", "ivan@example.com", "+79991234567"} {
		if strings.Contains(out, leak) {
			t.Errorf("leaked %q in %q", leak, out)
		}
	}
}

func TestMaskCSV_NotCSV(t *testing.T) {
	in := "просто текст, без структуры"
	if got := MaskCSV(in); got != in {
		t.Errorf("non-CSV modified: %q", got)
	}
}
