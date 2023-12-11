package server

import (
	"context"
	"errors"
	"net/http"

	"github.com/Schalure/gofermart/internal/gofermart"
)

//go:generate mockgen -destination=../mocks/mock_tokencheker.go -package=mocks github.com/Schalure/gofermart/internal/server TokenCheker
type TokenCheker interface {
	CheckValidJWT(tokenString string) (string, error)
}

type Middleware struct {
	logger      gofermart.Loggerer
	tokenCheker TokenCheker
}

func NewMidleware(logger gofermart.Loggerer, tokenCheker TokenCheker) *Middleware {

	return &Middleware{
		logger:      logger,
		tokenCheker: tokenCheker,
	}
}

type contextKey string

const contextLoginKey contextKey = "LoginKey"

func (m *Middleware) WithAuthentication(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		pc := "func (m *Middleware) WithAuthentication(h http.Handler) http.Handler"

		tokenCookie, err := r.Cookie(authorizationCookie)
		if err != nil {
			m.logger.Infow(
				pc,
				"message", "no authentication token found",
				"error", err,
			)
			http.Error(w, errors.New("no authentication token found").Error(), http.StatusUnauthorized)
			return
		}

		login, err := m.tokenCheker.CheckValidJWT(tokenCookie.Value)
		if err != nil {
			m.logger.Infow(
				pc,
				"message", "token failed",
				"error", err,
			)
			http.Error(w, errors.New("token failed").Error(), http.StatusUnauthorized)
			return
		}

		h.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), contextLoginKey, login)))
	})
}
