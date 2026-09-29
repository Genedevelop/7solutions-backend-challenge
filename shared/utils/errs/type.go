package errs

type Kind string

const (
	KindBadRequest   Kind = "bad_request"
	KindValidation   Kind = "validation"
	KindUnauthorized Kind = "unauthorized"
	KindForbidden    Kind = "forbidden"
	KindNotFound     Kind = "not_found"
	KindConflict     Kind = "conflict"
	KindInternal     Kind = "internal"
)

type Error struct {
	Kind    Kind
	Code    string
	Message string
	Fields  map[string]string
	Cause   error
	Status  int
}

type Response struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`
}
