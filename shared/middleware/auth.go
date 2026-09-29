package middleware

import (
	"strings"

	echov4 "github.com/labstack/echo/v4"

	"github.com/Genedevelop/7solutions-backend-challenge/shared/utils/errs"
)

const userIDKey = "userID"

type TokenParser interface {
	Parse(token string) (userID string, err error)
}

func JWTAuth(tokens TokenParser) echov4.MiddlewareFunc {
	return func(next echov4.HandlerFunc) echov4.HandlerFunc {
		return func(c echov4.Context) error {
			// Step 1: expect "Authorization: Bearer <token>"
			scheme, token, ok := strings.Cut(c.Request().Header.Get(echov4.HeaderAuthorization), " ")
			if !ok || !strings.EqualFold(scheme, "Bearer") || token == "" {
				c.Response().Header().Set(echov4.HeaderWWWAuthenticate, "Bearer")
				return errs.ErrMissingToken
			}

			// Step 2: verify signature, algorithm, issuer and expiry
			userID, err := tokens.Parse(token)
			if err != nil {
				c.Response().Header().Set(echov4.HeaderWWWAuthenticate, `Bearer error="invalid_token"`)
				return errs.ErrInvalidToken
			}

			// Step 3: keep the caller id for the handlers
			c.Set(userIDKey, userID)
			return next(c)
		}
	}
}

func UserID(c echov4.Context) string {
	id, _ := c.Get(userIDKey).(string)
	return id
}
