package echo

import (
	"net/http"

	echov4 "github.com/labstack/echo/v4"

	"github.com/Genedevelop/7solutions-backend-challenge/shared/middleware"
)

func (h *authHandler) RegisterUser(c echov4.Context) error {
	var req registerRequest
	if err := middleware.Bind(c, &req); err != nil {
		return err
	}
	u, err := h.auth.Register(c.Request().Context(), req.Name, req.Email, req.Password)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, registerResponse{ID: u.ID, Name: u.Name, Email: u.Email, CreatedAt: u.CreatedAt})
}
