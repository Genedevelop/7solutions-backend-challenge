package logger

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/Genedevelop/7solutions-backend-challenge/shared/utils/errs"
)

func TestErrorAddsCodeAndKind(t *testing.T) {
	var buf bytes.Buffer
	Init(&buf, "info")

	Error("request failed", fmt.Errorf("wrap: %w", errs.Internal(errors.New("socket closed"), "find user")), "path", "/users")

	var entry map[string]any
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("not JSON: %q", buf.String())
	}
	want := map[string]any{
		"level":      "ERROR",
		"msg":        "request failed",
		"error":      "wrap: INTERNAL: find user: socket closed",
		"error_code": "INTERNAL",
		"error_kind": "internal",
		"path":       "/users",
	}
	for k, v := range want {
		if entry[k] != v {
			t.Errorf("%s = %v, want %v", k, entry[k], v)
		}
	}
}

func TestLevelFilter(t *testing.T) {
	var buf bytes.Buffer
	Init(&buf, "warn")
	Info("hidden")
	Debug("hidden")
	Warn("shown")
	if bytes.Contains(buf.Bytes(), []byte("hidden")) || !bytes.Contains(buf.Bytes(), []byte("shown")) {
		t.Errorf("level filter wrong: %s", buf.String())
	}
}
