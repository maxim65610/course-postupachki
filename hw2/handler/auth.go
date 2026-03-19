// Package handler содержит HTTP обработчики для всех эндпоинтов API.
// Каждый обработчик принимает HTTP запросы, валидирует их и вызывает соответствующие сервисы.
package handler

import (
	"hw2/service"
	"net/http"
)

// AuthHandler обрабатывает запросы аутентификации.
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler создает новый экземпляр AuthHandler.
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type authRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type tokenResponse struct {
	Token string `json:"token"`
}

// Register обрабатывает POST /auth/register.
func (handler *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req authRequest
	if err := decodeJSON(r, &req); err != nil {
		writeHandlerError(w, err)
		return
	}

	token, err := handler.authService.Register(req.Username, req.Password)
	if err != nil {
		writeHandlerError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, tokenResponse{Token: token})
}

// Login обрабатывает POST /auth/login.
func (handler *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req authRequest
	if err := decodeJSON(r, &req); err != nil {
		writeHandlerError(w, err)
		return
	}
	token, err := handler.authService.Login(req.Username, req.Password)
	if err != nil {
		writeHandlerError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, tokenResponse{Token: token})
}
