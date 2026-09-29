package echo

import "time"

type updateUserRequest struct {
	Name  *string `json:"name"`
	Email *string `json:"email"`
}

type userResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}
