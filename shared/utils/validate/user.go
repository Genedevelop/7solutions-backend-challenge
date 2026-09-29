package validate

import (
	"fmt"
	"net/mail"
	"strings"
	"unicode/utf8"
)

const (
	maxNameLength     = 100
	minPasswordLength = 8
	maxPasswordBytes  = 72
)

func NormalizeName(name string) string {
	return strings.TrimSpace(name)
}

func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func Name(name string) string {
	switch {
	case name == "":
		return "is required"
	case utf8.RuneCountInString(name) > maxNameLength:
		return fmt.Sprintf("must be at most %d characters", maxNameLength)
	}
	return ""
}

func Email(email string) string {
	if email == "" {
		return "is required"
	}
	// Step 1: ParseAddress also accepts "Name <a@b.c>", so the parsed address must equal the input
	addr, err := mail.ParseAddress(email)
	if err != nil || addr.Address != email {
		return "must be a valid email address"
	}
	// Step 2: require a dotted domain, "a@localhost" is not useful for a real user
	domain := email[strings.LastIndex(email, "@")+1:]
	if !strings.Contains(domain, ".") || strings.HasSuffix(domain, ".") {
		return "must be a valid email address"
	}
	return ""
}

func Password(password string) string {
	// Step 1: cap at 72 bytes because bcrypt silently ignores the rest
	switch {
	case password == "":
		return "is required"
	case utf8.RuneCountInString(password) < minPasswordLength:
		return fmt.Sprintf("must be at least %d characters", minPasswordLength)
	case len(password) > maxPasswordBytes:
		return fmt.Sprintf("must be at most %d bytes", maxPasswordBytes)
	}
	return ""
}
