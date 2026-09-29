package domain

import (
	"context"

	"github.com/Genedevelop/7solutions-backend-challenge/internal/user/domain/entity"
)

func (s *UserService) Get(ctx context.Context, id string) (*entity.User, error) {
	return s.repo.GetByID(ctx, id)
}
