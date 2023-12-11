package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/Schalure/gofermart/internal/gofermart/gofermaterrors"
)

// User autotification type
type authenticationData struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// User registration handler. POST /api/user/register
func (h *Handler) UserRegistration(response http.ResponseWriter, request *http.Request) {

	authData, err := getAuthenticationData(request.Body)
	if err != nil {
		http.Error(response, "error reading user login and password", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(request.Context(), 5*time.Second)
	defer cancel()
	token, err := h.userManager.CreateUser(ctx, authData.Login, authData.Password)
	if err != nil {
		if errors.Is(err, gofermaterrors.LoginAlreadyTaken) {
			http.Error(response, gofermaterrors.LoginAlreadyTaken.Error(), http.StatusConflict)
			return
		}
		http.Error(response, gofermaterrors.Internal.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(response, &http.Cookie{
		Name:  authorizationCookie,
		Value: token,
	})
	response.WriteHeader(http.StatusOK)
}

// User authentication handler. POST /api/user/login
func (h *Handler) UserAuthentication(response http.ResponseWriter, request *http.Request) {

	authData, err := getAuthenticationData(request.Body)
	if err != nil {
		http.Error(response, "error reading user login and password", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(request.Context(), 5*time.Second)
	defer cancel()
	token, err := h.userManager.AuthenticationUser(ctx, authData.Login, authData.Password)
	if err != nil {
		if errors.Is(err, gofermaterrors.InvalidLoginPassword) {
			http.Error(response, err.Error(), http.StatusUnauthorized)
			return
		}
		http.Error(response, gofermaterrors.Internal.Error(), http.StatusInternalServerError)
	}

	http.SetCookie(response, &http.Cookie{
		Name:  authorizationCookie,
		Value: token,
	})
	response.WriteHeader(http.StatusOK)
}

// Get authentication data
func getAuthenticationData(r io.Reader) (authenticationData, error) {

	var buf bytes.Buffer
	var data authenticationData

	// читаем тело запроса
	if _, err := buf.ReadFrom(r); err != nil {
		return data, err
	}

	err := json.Unmarshal(buf.Bytes(), &data)

	return data, err
}
