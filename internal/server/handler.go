package server

import (
	"context"
)

const (
	authorizationCookie = "Authorization"
)

// Interface for interaction with users
//
//go:generate mockgen -destination=../mocks/mock_usermanager.go -package=mocks github.com/Schalure/gofermart/internal/server UserManager
type UserManager interface {
	CreateUser(ctx context.Context, login, password string) (string, error)
	AuthenticationUser(ctx context.Context, login, password string) (string, error)
}

// Main handler object struct
type Handler struct {
	userManager UserManager
}

// Constructor for Handler type
func NewHandler(userManager UserManager) *Handler {

	return &Handler{
		userManager: userManager,
	}
}
