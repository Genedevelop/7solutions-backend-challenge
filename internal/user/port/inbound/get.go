package inbound

import (
	"context"

	"github.com/Genedevelop/7solutions-backend-challenge/internal/user/domain/entity"
)

type IGetUserUseCase interface {
	Get(ctx context.Context, id string) (*entity.User, error)
}
