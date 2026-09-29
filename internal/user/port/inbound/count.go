package inbound

import "context"

type ICountUsersUseCase interface {
	Count(ctx context.Context) (int64, error)
}
