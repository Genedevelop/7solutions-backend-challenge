package inbound

import (
	"context"

	"github.com/Genedevelop/7solutions-backend-challenge/internal/user/domain/entity"
)

type IUpdateUserUseCase interface {
	Update(ctx context.Context, actorID, id string, name, email *string) (*entity.User, error)
}
