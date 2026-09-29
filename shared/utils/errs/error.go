package errs

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

func New(kind Kind, code, message string) *Error {
	return &Error{Kind: kind, Code: code, Message: message}
}

func Validation(fields map[string]string) *Error {
	return ErrValidation.WithFields(fields)
}

func Internal(cause error, op string) *Error {
	return ErrInternal.WithCause(fmt.Errorf("%s: %w", op, cause))
}

func (e *Error) WithMessage(message string) *Error {
	c := *e
	c.Message = message
	return &c
}

func (e *Error) WithFields(fields map[string]string) *Error {
	c := *e
	c.Fields = fields
	return &c
}

func (e *Error) WithCause(cause error) *Error {
	c := *e
	c.Cause = cause
	return &c
}

func (e *Error) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Code, e.Cause)
	}
	if len(e.Fields) > 0 {
		return fmt.Sprintf("%s: %s %v", e.Code, e.Message, e.Fields)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *Error) Unwrap() error {
	return e.Cause
}

func (e *Error) Is(target error) bool {
	// Step 1: match by code, so errors.Is(err, ErrUserNotFound) still holds for a copy carrying extra detail
	t, ok := target.(*Error)
	return ok && t.Code == e.Code
}

func From(err error) *Error {
	// Step 1: the client hung up, this is not a server failure even when a driver error wraps it
	if errors.Is(err, context.Canceled) {
		return ErrClientClosed.WithCause(err)
	}
	// Step 2: already a central error, use it as is
	var e *Error
	if errors.As(err, &e) {
		return e
	}
	// Step 3: anything unknown is treated as internal and keeps the original as cause
	return ErrInternal.WithCause(err)
}

func (e *Error) HTTPStatus() int {
	if e.Status != 0 {
		return e.Status
	}
	switch e.Kind {
	case KindBadRequest, KindValidation:
		return http.StatusBadRequest
	case KindUnauthorized:
		return http.StatusUnauthorized
	case KindForbidden:
		return http.StatusForbidden
	case KindNotFound:
		return http.StatusNotFound
	case KindConflict:
		return http.StatusConflict
	}
	return http.StatusInternalServerError
}

func (e *Error) Response() Response {
	return Response{Code: e.Code, Message: e.Message, Fields: e.Fields}
}
