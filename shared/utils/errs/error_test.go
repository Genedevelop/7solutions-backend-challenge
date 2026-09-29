package errs

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"
)

var errNotFound = New(KindNotFound, "USER_NOT_FOUND", "user not found")

func TestIsMatchesByCode(t *testing.T) {
	wrapped := fmt.Errorf("get user: %w", errNotFound)
	if !errors.Is(wrapped, errNotFound) {
		t.Error("wrapped sentinel not matched")
	}
	if errors.Is(New(KindNotFound, "OTHER", "x"), errNotFound) {
		t.Error("different code matched")
	}
}

func TestFrom(t *testing.T) {
	cause := errors.New("connection refused")

	got := From(fmt.Errorf("list: %w", errNotFound))
	if got.Code != "USER_NOT_FOUND" || got.HTTPStatus() != http.StatusNotFound {
		t.Errorf("got %+v", got)
	}

	got = From(cause)
	if got.Kind != KindInternal || got.HTTPStatus() != http.StatusInternalServerError || !errors.Is(got, cause) {
		t.Errorf("unknown error not mapped to internal: %+v", got)
	}
	if got.Response().Message != "internal server error" {
		t.Errorf("internal detail leaked: %+v", got.Response())
	}
}

func TestInternalKeepsCause(t *testing.T) {
	cause := errors.New("disk full")
	e := Internal(cause, "insert user")
	if !errors.Is(e, cause) || e.Error() != "INTERNAL: insert user: disk full" {
		t.Errorf("got %q", e.Error())
	}
}

func TestHTTPStatus(t *testing.T) {
	tests := map[Kind]int{
		KindBadRequest:   400,
		KindValidation:   400,
		KindUnauthorized: 401,
		KindForbidden:    403,
		KindNotFound:     404,
		KindConflict:     409,
		KindInternal:     500,
	}
	for kind, want := range tests {
		if got := (&Error{Kind: kind}).HTTPStatus(); got != want {
			t.Errorf("%s: got %d, want %d", kind, got, want)
		}
	}
}

func TestFromClientClosed(t *testing.T) {
	got := From(Internal(context.Canceled, "find user"))
	if got.Code != "CLIENT_CLOSED_REQUEST" || got.HTTPStatus() != 499 {
		t.Errorf("cancelled request mapped to %+v", got)
	}
}
