package echo

import (
	echov4 "github.com/labstack/echo/v4"

	"github.com/Genedevelop/7solutions-backend-challenge/internal/user/domain/entity"
	"github.com/Genedevelop/7solutions-backend-challenge/internal/user/port/inbound"
)

type userHandler struct {
	users inbound.IUserUseCase
}

func NewUserHandler(users inbound.IUserUseCase) *userHandler {
	return &userHandler{users: users}
}

func (h *userHandler) Register(e *echov4.Echo, auth echov4.MiddlewareFunc) {
	// Step 1: auth goes on each route, a group-level middleware makes Echo answer 404 instead of 405 for unknown methods
	e.GET("/users", h.ListUsers, auth)
	e.GET("/users/:id", h.GetUser, auth)
	e.PATCH("/users/:id", h.UpdateUser, auth)
	e.DELETE("/users/:id", h.DeleteUser, auth)
}

func toUserResponse(u *entity.User) userResponse {
	return userResponse{ID: u.ID, Name: u.Name, Email: u.Email, CreatedAt: u.CreatedAt}
}
