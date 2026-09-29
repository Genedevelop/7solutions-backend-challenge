package mock

import (
	"context"
	"fmt"

	"github.com/Genedevelop/7solutions-backend-challenge/internal/authentication/port/outbound"
	"github.com/Genedevelop/7solutions-backend-challenge/shared/utils/errs"
)

var _ outbound.ICredentialWriter = (*CredentialStore)(nil)

func (s *CredentialStore) Create(_ context.Context, account *outbound.Account) error {
	if s.Err != nil {
		return s.Err
	}
	if _, taken := s.ByEmail[account.Email]; taken {
		return errs.ErrEmailTaken
	}
	account.ID = fmt.Sprintf("user-%d", len(s.ByEmail)+1)
	s.ByEmail[account.Email] = *account
	return nil
}
