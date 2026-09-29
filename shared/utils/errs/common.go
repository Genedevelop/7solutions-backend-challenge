package errs

import "net/http"

var (
	ErrInternal         = New(KindInternal, "INTERNAL", "internal server error")
	ErrValidation       = New(KindValidation, "VALIDATION_FAILED", "validation failed")
	ErrConfigInvalid    = New(KindInternal, "CONFIG_INVALID", "invalid configuration")
	ErrBadRequest       = New(KindBadRequest, "BAD_REQUEST", "bad request")
	ErrInvalidJSON      = New(KindBadRequest, "INVALID_JSON", "request body must be valid JSON")
	ErrRouteNotFound    = &Error{Kind: KindNotFound, Code: "ROUTE_NOT_FOUND", Message: "route not found"}
	ErrMethodNotAllowed = &Error{Kind: KindBadRequest, Code: "METHOD_NOT_ALLOWED", Message: "method not allowed", Status: http.StatusMethodNotAllowed}
	ErrBodyTooLarge     = &Error{Kind: KindBadRequest, Code: "BODY_TOO_LARGE", Message: "request body too large", Status: http.StatusRequestEntityTooLarge}
	ErrClientClosed     = &Error{Kind: KindBadRequest, Code: "CLIENT_CLOSED_REQUEST", Message: "client closed request", Status: 499}
)
