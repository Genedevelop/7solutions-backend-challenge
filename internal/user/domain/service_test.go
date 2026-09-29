package domain_test

import (
	"context"
	"testing"

	"github.com/Genedevelop/7solutions-backend-challenge/internal/user/domain"
	"github.com/Genedevelop/7solutions-backend-challenge/internal/user/domain/entity"
	"github.com/Genedevelop/7solutions-backend-challenge/internal/user/port/outbound/mock"
)

func setup(t *testing.T) (*domain.UserService, *mock.UserRepository) {
	t.Helper()
	repo := mock.NewUserRepository()
	return domain.NewUserService(repo), repo
}

func seed(t *testing.T, repo *mock.UserRepository, email string) *entity.User {
	t.Helper()
	u := &entity.User{Name: "Test", Email: email}
	if err := repo.Create(context.Background(), u); err != nil {
		t.Fatalf("seed %s: %v", email, err)
	}
	return u
}

func ptr(s string) *string { return &s }
