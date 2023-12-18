package server

import (
	"context"

	"github.com/Schalure/gofermart/internal/storage"
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
	GetUserInfo(ctx context.Context, login string) (storage.User, error)
}

//go:generate mockgen -destination=../mocks/mock_ordermanager.go -package=mocks github.com/Schalure/gofermart/internal/server OrderManager
type OrderManager interface {
	LoadOrder(ctx context.Context, login, orderNumber string) error
	GetOrders(ctx context.Context, login string) ([]storage.Order, error)
}

//go:generate mockgen -destination=../mocks/mock_loyaltysystemmanager.go -package=mocks github.com/Schalure/gofermart/internal/server LoyaltySystemManager
type LoyaltySystemManager interface {
	Withdraw(ctx context.Context, login, orderNumber string, sum float64) error
	GetWithdraws(ctx context.Context, login string) ([]storage.Order, error)
}

// Main handler object struct
type Handler struct {
	userManager          UserManager
	orderManager         OrderManager
	loyaltySystemManager LoyaltySystemManager
}

// Constructor for Handler type
func NewHandler(userManager UserManager, orderManager OrderManager, loyaltySystemManager LoyaltySystemManager) *Handler {

	return &Handler{
		userManager:          userManager,
		orderManager:         orderManager,
		loyaltySystemManager: loyaltySystemManager,
	}
}

//	Get login from request context
func (h *Handler) getLoginFromContext(ctx context.Context) string {

	login := ctx.Value(contextLoginKey)
	l, ok := login.(string)
	if !ok {
		return ""
	}
	return l
}
