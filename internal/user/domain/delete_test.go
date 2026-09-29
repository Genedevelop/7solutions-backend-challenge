package domain_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Genedevelop/7solutions-backend-challenge/shared/utils/errs"
)

func TestDelete(t *testing.T) {
	ctx := context.Background()
	svc, repo := setup(t)
	a := seed(t, repo, "a@example.com")
	b := seed(t, repo, "b@example.com")

	if err := svc.Delete(ctx, a.ID, b.ID); !errors.Is(err, errs.ErrForbidden) {
		t.Errorf("other user: want ErrForbidden, got %v", err)
	}
	if err := svc.Delete(ctx, a.ID, a.ID); err != nil {
		t.Fatal(err)
	}
	if n, _ := svc.Count(ctx); n != 1 {
		t.Errorf("count = %d, want 1", n)
	}
}
