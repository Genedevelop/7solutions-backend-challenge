package outbound

import "time"

type Credential struct {
	UserID       string
	PasswordHash string
}

type Account struct {
	ID           string
	Name         string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}
