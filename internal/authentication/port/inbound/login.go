package inbound

import (
	"context"
)

type ILoginUseCase interface {
	Login(ctx context.Context, email, password string) (*LoginResult, error)
}
