package entity

import (
	"time"

	"github.com/Genedevelop/7solutions-backend-challenge/shared/utils/errs"
	"github.com/Genedevelop/7solutions-backend-challenge/shared/utils/validate"
)

type User struct {
	ID        string
	Name      string
	Email     string
	CreatedAt time.Time
}

func (u *User) ApplyUpdate(name, email *string) error {
	// Step 1: an empty body is a client mistake, not a no-op
	if name == nil && email == nil {
		return errs.Validation(map[string]string{"body": "name or email is required"})
	}

	// Step 2: validate the new values before touching the user
	fields := map[string]string{}
	var newName, newEmail string
	if name != nil {
		newName = validate.NormalizeName(*name)
		if msg := validate.Name(newName); msg != "" {
			fields["name"] = msg
		}
	}
	if email != nil {
		newEmail = validate.NormalizeEmail(*email)
		if msg := validate.Email(newEmail); msg != "" {
			fields["email"] = msg
		}
	}
	if len(fields) > 0 {
		return errs.Validation(fields)
	}

	// Step 3: apply only the fields that were sent
	if name != nil {
		u.Name = newName
	}
	if email != nil {
		u.Email = newEmail
	}
	return nil
}
