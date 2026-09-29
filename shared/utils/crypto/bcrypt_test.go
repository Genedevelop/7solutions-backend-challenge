package crypto

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestBcrypt(t *testing.T) {
	h := NewBcryptHasher(bcrypt.MinCost)
	hash, err := h.Hash("password123")
	if err != nil {
		t.Fatal(err)
	}
	if hash == "password123" {
		t.Fatal("password stored in plain text")
	}
	if err := h.Compare(hash, "password123"); err != nil {
		t.Errorf("correct password rejected: %v", err)
	}
	if err := h.Compare(hash, "wrong"); err == nil {
		t.Error("wrong password accepted")
	}
}
