package echo

import (
	"net/http"

	echov4 "github.com/labstack/echo/v4"
)

func (h *userHandler) GetUser(c echov4.Context) error {
	u, err := h.users.Get(c.Request().Context(), c.Param("id"))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, toUserResponse(u))
}
