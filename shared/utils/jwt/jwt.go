package jwt

import (
	"fmt"
	"time"

	gojwt "github.com/golang-jwt/jwt/v5"

	"github.com/Genedevelop/7solutions-backend-challenge/shared/utils/errs"
)

const minSecretLength = 32

type JWTManager struct {
	secret []byte
	ttl    time.Duration
	issuer string
	now    func() time.Time
}

func NewJWTManager(secret string, ttl time.Duration, issuer string) (*JWTManager, error) {
	// Step 1: HS256 is only as strong as the secret, refuse short ones at startup
	if len(secret) < minSecretLength {
		return nil, errs.ErrConfigInvalid.WithMessage(fmt.Sprintf("jwt secret must be at least %d bytes", minSecretLength))
	}
	return &JWTManager{secret: []byte(secret), ttl: ttl, issuer: issuer, now: time.Now}, nil
}

func (m *JWTManager) Generate(userID string) (string, time.Time, error) {
	// Step 1: put the user id in "sub" with issue and expiry times, in UTC seconds since that is all a JWT can hold
	now := m.now().UTC().Truncate(time.Second)
	expiresAt := now.Add(m.ttl)
	claims := gojwt.RegisteredClaims{
		Subject:   userID,
		Issuer:    m.issuer,
		IssuedAt:  gojwt.NewNumericDate(now),
		ExpiresAt: gojwt.NewNumericDate(expiresAt),
	}

	// Step 2: sign with HMAC-SHA256
	token, err := gojwt.NewWithClaims(gojwt.SigningMethodHS256, claims).SignedString(m.secret)
	if err != nil {
		return "", time.Time{}, errs.Internal(err, "sign token")
	}
	return token, expiresAt, nil
}

func (m *JWTManager) Parse(tokenString string) (string, error) {
	// Step 1: pin HS256 so "alg: none" or a swapped algorithm is rejected, and require exp plus our issuer
	var claims gojwt.RegisteredClaims
	_, err := gojwt.ParseWithClaims(tokenString, &claims,
		func(*gojwt.Token) (any, error) { return m.secret, nil },
		gojwt.WithValidMethods([]string{gojwt.SigningMethodHS256.Alg()}),
		gojwt.WithIssuer(m.issuer),
		gojwt.WithExpirationRequired(),
		gojwt.WithTimeFunc(m.now),
	)
	if err != nil {
		return "", errs.ErrInvalidToken.WithCause(err)
	}

	// Step 2: a token without a subject cannot be tied to a user
	if claims.Subject == "" {
		return "", errs.ErrInvalidToken
	}
	return claims.Subject, nil
}
