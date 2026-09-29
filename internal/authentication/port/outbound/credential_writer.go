package outbound

import "context"

type ICredentialWriter interface {
	Create(ctx context.Context, account *Account) error
}
