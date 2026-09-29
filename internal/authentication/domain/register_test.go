package domain_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Genedevelop/7solutions-backend-challenge/internal/authentication/domain"
	"github.com/Genedevelop/7solutions-backend-challenge/internal/authentication/port/outbound/mock"
	"github.com/Genedevelop/7solutions-backend-challenge/shared/utils/errs"
)

func TestRegister(t *testing.T) {
	ctx := context.Background()
	store := mock.NewCredentialStore()
	svc, err := domain.NewAuthService(store, store, mock.PlainHasher{}, mock.Tokens{})
	if err != nil {
		t.Fatal(err)
	}

	t.Run("normalizes, hashes and returns the new user", func(t *testing.T) {
		u, err := svc.Register(ctx, "  Alice ", " Alice@Example.COM ", "password123")
		if err != nil {
			t.Fatal(err)
		}
		if u.ID == "" || u.Name != "Alice" || u.Email != "alice@example.com" || u.CreatedAt.IsZero() {
			t.Errorf("got %+v", u)
		}
		if got := store.ByEmail["alice@example.com"].PasswordHash; got != "hashed:password123" {
			t.Errorf("stored hash = %q", got)
		}
	})

	t.Run("duplicate email in any case", func(t *testing.T) {
		if _, err := svc.Register(ctx, "Other", "ALICE@example.com", "password123"); !errors.Is(err, errs.ErrEmailTaken) {
			t.Fatalf("want ErrEmailTaken, got %v", err)
		}
	})

	t.Run("every invalid field reported at once", func(t *testing.T) {
		_, err := svc.Register(ctx, "", "bad", "short")
		var e *errs.Error
		if !errors.As(err, &e) || e.Kind != errs.KindValidation || len(e.Fields) != 3 {
			t.Fatalf("want 3 field errors, got %v", err)
		}
		if len(store.ByEmail) != 1 {
			t.Errorf("invalid user was stored")
		}
	})

	t.Run("registered user can log in", func(t *testing.T) {
		res, err := svc.Login(ctx, "alice@example.com", "password123")
		if err != nil || res.Token == "" {
			t.Fatalf("Login = %+v, %v", res, err)
		}
	})
}
