package config

import (
	"testing"
	"time"
)

func TestLoadDefaults(t *testing.T) {
	t.Setenv("JWT_SECRET", "0123456789abcdef0123456789abcdef")
	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.HTTPAddr != ":8080" || cfg.JWTTTL != time.Hour || cfg.UserCountInterval != 10*time.Second {
		t.Errorf("got %+v", cfg)
	}
}

func TestLoadErrors(t *testing.T) {
	t.Setenv("JWT_SECRET", "")
	if _, err := Load(); err == nil {
		t.Error("missing JWT_SECRET should fail")
	}

	t.Setenv("JWT_SECRET", "0123456789abcdef0123456789abcdef")
	t.Setenv("JWT_TTL", "one hour")
	if _, err := Load(); err == nil {
		t.Error("bad JWT_TTL should fail")
	}
}
