package outbound

import (
	"context"

	"github.com/Genedevelop/7solutions-backend-challenge/internal/user/domain/entity"
)

type IUserRepository interface {
	GetByID(ctx context.Context, id string) (*entity.User, error)
	List(ctx context.Context) ([]*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context) (int64, error)
}
