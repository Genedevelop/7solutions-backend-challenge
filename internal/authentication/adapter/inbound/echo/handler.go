package echo

import (
	echov4 "github.com/labstack/echo/v4"

	"github.com/Genedevelop/7solutions-backend-challenge/internal/authentication/port/inbound"
)

type authHandler struct {
	auth inbound.IAuthUseCase
}

func NewAuthHandler(auth inbound.IAuthUseCase) *authHandler {
	return &authHandler{auth: auth}
}

func (h *authHandler) Register(e *echov4.Echo) {
	e.POST("/auth/register", h.RegisterUser)
	e.POST("/auth/login", h.Login)
}
