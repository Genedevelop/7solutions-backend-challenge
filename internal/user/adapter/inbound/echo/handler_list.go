package echo

import (
	"net/http"

	echov4 "github.com/labstack/echo/v4"
)

func (h *userHandler) ListUsers(c echov4.Context) error {
	users, err := h.users.List(c.Request().Context())
	if err != nil {
		return err
	}
	out := make([]userResponse, 0, len(users))
	for _, u := range users {
		out = append(out, toUserResponse(u))
	}
	return c.JSON(http.StatusOK, out)
}
