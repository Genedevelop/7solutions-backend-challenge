package config

import (
	"testing"
	"time"
)

func clearEnv(t *testing.T) {
	t.Helper()
	for _, key := range []string{
		"HTTP_ADDR", "LOG_LEVEL", "SHUTDOWN_TIMEOUT", "MONGO_URI", "MONGO_DB",
		"JWT_SECRET", "JWT_ISSUER", "JWT_TTL", "USER_COUNT_INTERVAL",
	} {
		t.Setenv(key, "")
	}
}

func TestLoadDefaults(t *testing.T) {
	clearEnv(t)
	t.Setenv("JWT_SECRET", "0123456789abcdef0123456789abcdef")
	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.HTTPAddr != ":8080" || cfg.MongoDB != "user_api" || cfg.JWTTTL != time.Hour || cfg.UserCountInterval != 10*time.Second {
		t.Errorf("unexpected defaults: addr=%s db=%s ttl=%s interval=%s", cfg.HTTPAddr, cfg.MongoDB, cfg.JWTTTL, cfg.UserCountInterval)
	}
}

func TestLoadErrors(t *testing.T) {
	clearEnv(t)
	if _, err := Load(); err == nil {
		t.Error("missing JWT_SECRET should fail")
	}

	t.Setenv("JWT_SECRET", "0123456789abcdef0123456789abcdef")
	t.Setenv("JWT_TTL", "one hour")
	if _, err := Load(); err == nil {
		t.Error("bad JWT_TTL should fail")
	}
}
