package echo

import (
	"net/http"

	echov4 "github.com/labstack/echo/v4"

	"github.com/Genedevelop/7solutions-backend-challenge/shared/middleware"
)

func (h *userHandler) UpdateUser(c echov4.Context) error {
	var req updateUserRequest
	if err := middleware.Bind(c, &req); err != nil {
		return err
	}
	u, err := h.users.Update(c.Request().Context(), middleware.UserID(c), c.Param("id"), req.Name, req.Email)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, toUserResponse(u))
}
