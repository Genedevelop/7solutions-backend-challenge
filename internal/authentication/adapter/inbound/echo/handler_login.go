package echo

import (
	"net/http"

	echov4 "github.com/labstack/echo/v4"

	"github.com/Genedevelop/7solutions-backend-challenge/shared/middleware"
)

func (h *authHandler) Login(c echov4.Context) error {
	var req loginRequest
	if err := middleware.Bind(c, &req); err != nil {
		return err
	}
	res, err := h.auth.Login(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, loginResponse{AccessToken: res.Token, TokenType: "Bearer", ExpiresAt: res.ExpiresAt})
}
