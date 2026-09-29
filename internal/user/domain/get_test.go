package domain_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Genedevelop/7solutions-backend-challenge/shared/utils/errs"
)

func TestGetAndList(t *testing.T) {
	ctx := context.Background()
	svc, repo := setup(t)
	a := seed(t, repo, "a@example.com")
	seed(t, repo, "b@example.com")

	if got, err := svc.Get(ctx, a.ID); err != nil || got.Email != a.Email {
		t.Errorf("Get = %+v, %v", got, err)
	}
	if _, err := svc.Get(ctx, "missing"); !errors.Is(err, errs.ErrUserNotFound) {
		t.Errorf("want ErrUserNotFound, got %v", err)
	}
	if users, err := svc.List(ctx); err != nil || len(users) != 2 {
		t.Errorf("List = %d users, %v", len(users), err)
	}

	boom := errors.New("mongo down")
	repo.Err = boom
	if _, err := svc.List(ctx); !errors.Is(err, boom) {
		t.Errorf("List with db error: want %v, got %v", boom, err)
	}
}
