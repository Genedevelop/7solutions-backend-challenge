package domain_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Genedevelop/7solutions-backend-challenge/internal/authentication/domain"
	"github.com/Genedevelop/7solutions-backend-challenge/internal/authentication/port/outbound"
	"github.com/Genedevelop/7solutions-backend-challenge/internal/authentication/port/outbound/mock"
	"github.com/Genedevelop/7solutions-backend-challenge/shared/utils/errs"
)

func TestLogin(t *testing.T) {
	ctx := context.Background()
	reader := mock.NewCredentialStore()
	reader.ByEmail["alice@example.com"] = outbound.Account{ID: "user-1", Email: "alice@example.com", PasswordHash: "hashed:password123"}
	svc, err := domain.NewAuthService(reader, reader, mock.PlainHasher{}, mock.Tokens{})
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name     string
		email    string
		password string
		wantErr  error
	}{
		{"ok, email case and spaces ignored", " ALICE@example.com ", "password123", nil},
		{"wrong password", "alice@example.com", "nope-nope", errs.ErrInvalidCredentials},
		{"unknown email", "bob@example.com", "password123", errs.ErrInvalidCredentials},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res, err := svc.Login(ctx, tc.email, tc.password)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("want %v, got %v", tc.wantErr, err)
			}
			if tc.wantErr == nil && res.Token != "token:user-1" {
				t.Errorf("token = %q", res.Token)
			}
		})
	}

	t.Run("database error is not reported as bad credentials", func(t *testing.T) {
		boom := errors.New("mongo down")
		reader.Err = boom
		defer func() { reader.Err = nil }()
		if _, err := svc.Login(ctx, "alice@example.com", "password123"); !errors.Is(err, boom) {
			t.Errorf("want %v, got %v", boom, err)
		}
	})
}
