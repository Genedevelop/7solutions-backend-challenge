package jwt

import (
	"strings"
	"testing"
	"time"

	gojwt "github.com/golang-jwt/jwt/v5"
)

const testSecret = "0123456789abcdef0123456789abcdef"

func newManager(t *testing.T) *JWTManager {
	t.Helper()
	m, err := NewJWTManager(testSecret, time.Hour, "user-api")
	if err != nil {
		t.Fatal(err)
	}
	return m
}

func TestJWTRoundTrip(t *testing.T) {
	m := newManager(t)
	token, exp, err := m.Generate("user-1")
	if err != nil {
		t.Fatal(err)
	}
	if time.Until(exp) < 59*time.Minute {
		t.Errorf("expiry too early: %v", exp)
	}
	sub, err := m.Parse(token)
	if err != nil || sub != "user-1" {
		t.Fatalf("Parse = %q, %v", sub, err)
	}
}

func TestJWTRejects(t *testing.T) {
	m := newManager(t)
	good, _, _ := m.Generate("user-1")

	expired := newManager(t)
	expired.now = func() time.Time { return time.Now().Add(-2 * time.Hour) }
	expiredToken, _, _ := expired.Generate("user-1")

	other, _ := NewJWTManager(strings.Repeat("x", 32), time.Hour, "user-api")
	wrongSecret, _, _ := other.Generate("user-1")

	otherIssuer, _ := NewJWTManager(testSecret, time.Hour, "someone-else")
	wrongIssuer, _, _ := otherIssuer.Generate("user-1")

	hs512, _ := gojwt.NewWithClaims(gojwt.SigningMethodHS512, gojwt.RegisteredClaims{
		Subject: "user-1", Issuer: "user-api", ExpiresAt: gojwt.NewNumericDate(time.Now().Add(time.Hour)),
	}).SignedString([]byte(testSecret))

	none, _ := gojwt.NewWithClaims(gojwt.SigningMethodNone, gojwt.RegisteredClaims{
		Subject: "user-1", Issuer: "user-api", ExpiresAt: gojwt.NewNumericDate(time.Now().Add(time.Hour)),
	}).SignedString(gojwt.UnsafeAllowNoneSignatureType)

	noExp, _ := gojwt.NewWithClaims(gojwt.SigningMethodHS256, gojwt.RegisteredClaims{
		Subject: "user-1", Issuer: "user-api",
	}).SignedString([]byte(testSecret))

	tests := map[string]string{
		"expired":      expiredToken,
		"wrong secret": wrongSecret,
		"wrong issuer": wrongIssuer,
		"HS512":        hs512,
		"alg none":     none,
		"no exp":       noExp,
		"tampered":     good[:len(good)-2] + "xx",
		"garbage":      "not.a.token",
	}
	for name, token := range tests {
		t.Run(name, func(t *testing.T) {
			if _, err := m.Parse(token); err == nil {
				t.Error("expected token to be rejected")
			}
		})
	}
}

func TestJWTShortSecret(t *testing.T) {
	if _, err := NewJWTManager("short", time.Hour, "user-api"); err == nil {
		t.Error("expected short secret to be rejected")
	}
}
