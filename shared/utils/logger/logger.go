package logger

import (
	"errors"
	"io"
	"log/slog"
	"os"
	"strings"
	"sync/atomic"

	"github.com/Genedevelop/7solutions-backend-challenge/shared/utils/errs"
)

var current atomic.Pointer[slog.Logger]

func init() {
	current.Store(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
}

func Init(w io.Writer, level string) {
	// Step 1: map the configured level, unknown values fall back to info
	var lvl slog.Level
	switch strings.ToLower(level) {
	case "debug":
		lvl = slog.LevelDebug
	case "warn":
		lvl = slog.LevelWarn
	case "error":
		lvl = slog.LevelError
	default:
		lvl = slog.LevelInfo
	}

	// Step 2: JSON so log collectors can index every field
	current.Store(slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{Level: lvl})))
}

func Get() *slog.Logger {
	return current.Load()
}

func Debug(msg string, args ...any) { current.Load().Debug(msg, args...) }

func Info(msg string, args ...any) { current.Load().Info(msg, args...) }

func Warn(msg string, args ...any) { current.Load().Warn(msg, args...) }

func Error(msg string, err error, args ...any) {
	// Step 1: always log the full error text, including wrapped causes the client never sees
	attrs := append([]any{"error", err.Error()}, args...)

	// Step 2: central errors also carry their code and kind, which is what alerts and dashboards filter on
	var e *errs.Error
	if errors.As(err, &e) {
		attrs = append(attrs, "error_code", e.Code, "error_kind", string(e.Kind))
	}
	current.Load().Error(msg, attrs...)
}
