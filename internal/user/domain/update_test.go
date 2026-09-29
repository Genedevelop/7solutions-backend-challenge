package domain_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Genedevelop/7solutions-backend-challenge/shared/utils/errs"
)

func TestUpdate(t *testing.T) {
	ctx := context.Background()
	svc, repo := setup(t)
	a := seed(t, repo, "a@example.com")
	b := seed(t, repo, "b@example.com")

	got, err := svc.Update(ctx, a.ID, a.ID, ptr("Alice"), ptr("ALICE@example.com"))
	if err != nil || got.Name != "Alice" || got.Email != "alice@example.com" {
		t.Fatalf("Update = %+v, %v", got, err)
	}
	if _, err := svc.Update(ctx, a.ID, b.ID, ptr("hacked"), nil); !errors.Is(err, errs.ErrForbidden) {
		t.Errorf("other user: want ErrForbidden, got %v", err)
	}
	if _, err := svc.Update(ctx, b.ID, b.ID, nil, ptr("alice@example.com")); !errors.Is(err, errs.ErrEmailTaken) {
		t.Errorf("taken email: want ErrEmailTaken, got %v", err)
	}
	var verr *errs.Error
	if _, err := svc.Update(ctx, b.ID, b.ID, nil, nil); !errors.As(err, &verr) || verr.Kind != errs.KindValidation {
		t.Errorf("empty update: want ValidationError, got %v", err)
	}
}
