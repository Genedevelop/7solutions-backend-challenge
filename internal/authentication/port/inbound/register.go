package inbound

import "context"

type IRegisterUseCase interface {
	Register(ctx context.Context, name, email, password string) (*RegisteredUser, error)
}
