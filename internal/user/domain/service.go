package domain

import (
	"github.com/Genedevelop/7solutions-backend-challenge/internal/user/port/inbound"
	"github.com/Genedevelop/7solutions-backend-challenge/internal/user/port/outbound"
)

var _ inbound.IUserUseCase = (*UserService)(nil)

type UserService struct {
	repo outbound.IUserRepository
}

func NewUserService(repo outbound.IUserRepository) *UserService {
	return &UserService{repo: repo}
}
