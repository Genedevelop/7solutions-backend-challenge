package mock

import (
	"errors"

	"github.com/Genedevelop/7solutions-backend-challenge/internal/authentication/port/outbound"
)

var _ outbound.IPasswordHasher = PlainHasher{}

type PlainHasher struct{}

func (PlainHasher) Hash(password string) (string, error) {
	return "hashed:" + password, nil
}

func (PlainHasher) Compare(hash, password string) error {
	if hash != "hashed:"+password {
		return errors.New("password mismatch")
	}
	return nil
}
