package errs

var (
	ErrMissingToken       = New(KindUnauthorized, "MISSING_TOKEN", "missing bearer token")
	ErrInvalidToken       = New(KindUnauthorized, "INVALID_TOKEN", "invalid or expired token")
	ErrInvalidCredentials = New(KindUnauthorized, "INVALID_CREDENTIALS", "invalid email or password")
	ErrCredentialNotFound = New(KindNotFound, "CREDENTIAL_NOT_FOUND", "credential not found")
)
