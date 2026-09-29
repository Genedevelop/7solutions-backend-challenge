package validate

import (
	"strings"
	"testing"
)

func TestEmail(t *testing.T) {
	tests := map[string]bool{
		"alice@example.com":         true,
		"a.b+tag@sub.example.co.th": true,
		"Alice <alice@example.com>": false,
		"alice@localhost":           false,
		"alice@example.":            false,
		"alice.example.com":         false,
		"":                          false,
	}
	for in, valid := range tests {
		if got := Email(in) == ""; got != valid {
			t.Errorf("Email(%q) valid = %v, want %v", in, got, valid)
		}
	}
}

func TestName(t *testing.T) {
	if Name("Alice") != "" || Name("") == "" || Name(strings.Repeat("n", 101)) == "" {
		t.Error("name rules wrong")
	}
}

func TestPassword(t *testing.T) {
	tests := map[string]bool{
		"password123":           true,
		"":                      false,
		"short":                 false,
		strings.Repeat("a", 72): true,
		strings.Repeat("a", 73): false,
	}
	for in, valid := range tests {
		if got := Password(in) == ""; got != valid {
			t.Errorf("Password(len %d) valid = %v, want %v", len(in), got, valid)
		}
	}
}

func TestNormalizeEmail(t *testing.T) {
	if got := NormalizeEmail("  Alice@Example.COM "); got != "alice@example.com" {
		t.Errorf("NormalizeEmail = %q", got)
	}
}
