package errs

var (
	ErrUserNotFound = New(KindNotFound, "USER_NOT_FOUND", "user not found")
	ErrEmailTaken   = New(KindConflict, "EMAIL_TAKEN", "email already registered")
	ErrForbidden    = New(KindForbidden, "FORBIDDEN", "cannot modify another user")
)
