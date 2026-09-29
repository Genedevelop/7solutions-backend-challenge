package domain

import (
	"github.com/Genedevelop/7solutions-backend-challenge/internal/authentication/port/inbound"
	"github.com/Genedevelop/7solutions-backend-challenge/internal/authentication/port/outbound"
	"github.com/Genedevelop/7solutions-backend-challenge/shared/utils/errs"
)

var _ inbound.IAuthUseCase = (*AuthService)(nil)

type AuthService struct {
	reader    outbound.ICredentialReader
	writer    outbound.ICredentialWriter
	hasher    outbound.IPasswordHasher
	tokens    outbound.ITokenIssuer
	dummyHash string
}

func NewAuthService(reader outbound.ICredentialReader, writer outbound.ICredentialWriter, hasher outbound.IPasswordHasher, tokens outbound.ITokenIssuer) (*AuthService, error) {
	// Step 1: pre-hash a throwaway password so a login for an unknown email costs the same as a real one
	dummy, err := hasher.Hash("not-a-real-password")
	if err != nil {
		return nil, errs.Internal(err, "prepare dummy hash")
	}
	return &AuthService{reader: reader, writer: writer, hasher: hasher, tokens: tokens, dummyHash: dummy}, nil
}
