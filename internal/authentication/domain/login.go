package domain

import (
	"context"
	"errors"

	"github.com/Genedevelop/7solutions-backend-challenge/internal/authentication/port/inbound"
	"github.com/Genedevelop/7solutions-backend-challenge/shared/utils/errs"
	"github.com/Genedevelop/7solutions-backend-challenge/shared/utils/validate"
)

func (s *AuthService) Login(ctx context.Context, email, password string) (*inbound.LoginResult, error) {
	// Step 1: emails are stored lower-cased, look the account up the same way
	cred, err := s.reader.GetByEmail(ctx, validate.NormalizeEmail(email))
	if errors.Is(err, errs.ErrCredentialNotFound) {
		// Step 2: burn the same bcrypt time and give the same error, so callers cannot probe which emails exist
		_ = s.hasher.Compare(s.dummyHash, password)
		return nil, errs.ErrInvalidCredentials
	}
	if err != nil {
		return nil, err
	}

	// Step 3: check the password
	if err := s.hasher.Compare(cred.PasswordHash, password); err != nil {
		return nil, errs.ErrInvalidCredentials
	}

	// Step 4: issue a JWT with the user id as subject
	token, expiresAt, err := s.tokens.Generate(cred.UserID)
	if err != nil {
		return nil, errs.Internal(err, "issue token")
	}
	return &inbound.LoginResult{Token: token, ExpiresAt: expiresAt}, nil
}
