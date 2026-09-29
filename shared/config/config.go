package config

import (
	"fmt"
	"os"
	"time"

	"github.com/Genedevelop/7solutions-backend-challenge/shared/utils/errs"
)

type Config struct {
	LogLevel          string
	HTTPAddr          string
	MongoURI          string
	MongoDB           string
	JWTSecret         string
	JWTIssuer         string
	JWTTTL            time.Duration
	UserCountInterval time.Duration
	ShutdownTimeout   time.Duration
}

func Load() (*Config, error) {
	// Step 1: read env with sensible local defaults
	cfg := &Config{
		LogLevel:  getEnv("LOG_LEVEL", "info"),
		HTTPAddr:  getEnv("HTTP_ADDR", ":8080"),
		MongoURI:  getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDB:   getEnv("MONGO_DB", "user_api"),
		JWTSecret: os.Getenv("JWT_SECRET"),
		JWTIssuer: getEnv("JWT_ISSUER", "user-api"),
	}

	// Step 2: parse durations, a typo should stop startup instead of silently using a default
	var err error
	if cfg.JWTTTL, err = getDuration("JWT_TTL", time.Hour); err != nil {
		return nil, err
	}
	if cfg.UserCountInterval, err = getDuration("USER_COUNT_INTERVAL", 10*time.Second); err != nil {
		return nil, err
	}
	if cfg.ShutdownTimeout, err = getDuration("SHUTDOWN_TIMEOUT", 10*time.Second); err != nil {
		return nil, err
	}

	// Step 3: the signing secret has no default on purpose
	if cfg.JWTSecret == "" {
		return nil, errs.ErrConfigInvalid.WithMessage("JWT_SECRET is required")
	}
	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

func getDuration(key string, fallback time.Duration) (time.Duration, error) {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		return fallback, nil
	}
	d, err := time.ParseDuration(v)
	if err != nil || d <= 0 {
		return 0, errs.ErrConfigInvalid.WithMessage(fmt.Sprintf("%s must be a positive duration like 10s, got %q", key, v))
	}
	return d, nil
}
