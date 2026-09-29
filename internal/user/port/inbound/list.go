package inbound

import (
	"context"

	"github.com/Genedevelop/7solutions-backend-challenge/internal/user/domain/entity"
)

type IListUsersUseCase interface {
	List(ctx context.Context) ([]*entity.User, error)
}
