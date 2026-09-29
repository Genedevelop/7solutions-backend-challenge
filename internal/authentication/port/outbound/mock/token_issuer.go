package mock

import (
	"time"

	"github.com/Genedevelop/7solutions-backend-challenge/internal/authentication/port/outbound"
)

var _ outbound.ITokenIssuer = Tokens{}

type Tokens struct{}

func (Tokens) Generate(userID string) (string, time.Time, error) {
	return "token:" + userID, time.Now().Add(time.Hour), nil
}
