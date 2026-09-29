package domain

import (
	"context"
	"time"

	"github.com/Genedevelop/7solutions-backend-challenge/internal/authentication/port/inbound"
	"github.com/Genedevelop/7solutions-backend-challenge/internal/authentication/port/outbound"
	"github.com/Genedevelop/7solutions-backend-challenge/shared/utils/errs"
	"github.com/Genedevelop/7solutions-backend-challenge/shared/utils/validate"
)

func (s *AuthService) Register(ctx context.Context, name, email, password string) (*inbound.RegisteredUser, error) {
	// Step 1: normalize so "Alice@X.com " and "alice@x.com" are the same account
	account := &outbound.Account{
		Name:      validate.NormalizeName(name),
		Email:     validate.NormalizeEmail(email),
		CreatedAt: time.Now().UTC().Truncate(time.Millisecond),
	}

	// Step 2: collect every invalid field at once instead of failing on the first one
	fields := map[string]string{}
	if msg := validate.Name(account.Name); msg != "" {
		fields["name"] = msg
	}
	if msg := validate.Email(account.Email); msg != "" {
		fields["email"] = msg
	}
	if msg := validate.Password(password); msg != "" {
		fields["password"] = msg
	}
	if len(fields) > 0 {
		return nil, errs.Validation(fields)
	}

	// Step 3: never store the plain password
	hash, err := s.hasher.Hash(password)
	if err != nil {
		return nil, errs.Internal(err, "hash password")
	}
	account.PasswordHash = hash

	// Step 4: save, the unique email index turns a duplicate into ErrEmailTaken
	if err := s.writer.Create(ctx, account); err != nil {
		return nil, err
	}
	return &inbound.RegisteredUser{ID: account.ID, Name: account.Name, Email: account.Email, CreatedAt: account.CreatedAt}, nil
}
