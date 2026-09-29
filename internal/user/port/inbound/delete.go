package inbound

import "context"

type IDeleteUserUseCase interface {
	Delete(ctx context.Context, actorID, id string) error
}
