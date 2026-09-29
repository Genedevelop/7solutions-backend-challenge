package outbound

import "time"

type ITokenIssuer interface {
	Generate(userID string) (token string, expiresAt time.Time, err error)
}
