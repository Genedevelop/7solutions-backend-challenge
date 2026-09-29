package inbound

import "time"

type LoginResult struct {
	Token     string
	ExpiresAt time.Time
}

type RegisteredUser struct {
	ID        string
	Name      string
	Email     string
	CreatedAt time.Time
}
