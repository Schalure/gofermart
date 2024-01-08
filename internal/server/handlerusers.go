package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/Schalure/gofermart/internal/gofermart"
)

// User autotification type
type authenticationData struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// User registration handler. POST /api/user/register
func (h *Handler) UserRegistration(response http.ResponseWriter, request *http.Request) {

	var buf bytes.Buffer
	var authData authenticationData

	if _, err := buf.ReadFrom(request.Body); err != nil {
		http.Error(response, "error reading user login and password", http.StatusBadRequest)
		return
	}

	err := json.Unmarshal(buf.Bytes(), &authData)
	if err != nil {
		http.Error(response, "error reading user login and password", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(request.Context(), 5*time.Second)
	defer cancel()
	token, err := h.userManager.CreateUser(ctx, authData.Login, authData.Password)
	if err != nil {
		if errors.Is(err, gofermart.ErrLoginAlreadyTaken) {
			http.Error(response, gofermart.ErrLoginAlreadyTaken.Error(), http.StatusConflict)
			return
		}
		http.Error(response, gofermart.ErrInternal.Error(), http.StatusInternalServerError)
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

	var buf bytes.Buffer
	var authData authenticationData

	if _, err := buf.ReadFrom(request.Body); err != nil {
		http.Error(response, "error reading user login and password", http.StatusBadRequest)
		return
	}

	err := json.Unmarshal(buf.Bytes(), &authData)
	if err != nil {
		http.Error(response, "error reading user login and password", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(request.Context(), 5*time.Second)
	defer cancel()
	token, err := h.userManager.AuthenticationUser(ctx, authData.Login, authData.Password)
	if err != nil {
		if errors.Is(err, gofermart.ErrInvalidLoginPassword) {
			http.Error(response, err.Error(), http.StatusUnauthorized)
			return
		}
		http.Error(response, gofermart.ErrInternal.Error(), http.StatusInternalServerError)
	}

	http.SetCookie(response, &http.Cookie{
		Name:  authorizationCookie,
		Value: token,
	})
	response.WriteHeader(http.StatusOK)
}
