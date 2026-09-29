package outbound

import (
	"context"
)

type ICredentialReader interface {
	GetByEmail(ctx context.Context, email string) (*Credential, error)
}
