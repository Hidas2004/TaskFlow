package v1services

import (
	"github.com/Hidas2004/TaskFlow/internal/dto"
)

type AuthService interface {
	Login(req *dto.LoginRequest) (*dto.AuthResponse, error)
	Register(req *dto.RegisterRequest) (*dto.AuthResponse, error)
}
