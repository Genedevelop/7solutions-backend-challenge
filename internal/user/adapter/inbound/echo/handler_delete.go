package echo

import (
	"net/http"

	echov4 "github.com/labstack/echo/v4"

	"github.com/Genedevelop/7solutions-backend-challenge/shared/middleware"
)

func (h *userHandler) DeleteUser(c echov4.Context) error {
	if err := h.users.Delete(c.Request().Context(), middleware.UserID(c), c.Param("id")); err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}
