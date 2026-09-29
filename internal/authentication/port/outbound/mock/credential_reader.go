package mock

import (
	"context"

	"github.com/Genedevelop/7solutions-backend-challenge/internal/authentication/port/outbound"
	"github.com/Genedevelop/7solutions-backend-challenge/shared/utils/errs"
)

var _ outbound.ICredentialReader = (*CredentialStore)(nil)

type CredentialStore struct {
	ByEmail map[string]outbound.Account
	Err     error
}

func NewCredentialStore() *CredentialStore {
	return &CredentialStore{ByEmail: map[string]outbound.Account{}}
}

func (s *CredentialStore) GetByEmail(_ context.Context, email string) (*outbound.Credential, error) {
	if s.Err != nil {
		return nil, s.Err
	}
	a, ok := s.ByEmail[email]
	if !ok {
		return nil, errs.ErrCredentialNotFound
	}
	return &outbound.Credential{UserID: a.ID, PasswordHash: a.PasswordHash}, nil
}
