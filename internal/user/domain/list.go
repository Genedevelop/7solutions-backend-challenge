package domain

import (
	"context"

	"github.com/Genedevelop/7solutions-backend-challenge/internal/user/domain/entity"
)

func (s *UserService) List(ctx context.Context) ([]*entity.User, error) {
	return s.repo.List(ctx)
}
