// Package service содержит бизнес-логику приложения
package service

import (
	"errors"
	"hw2/auth"
	"hw2/domain"
	"hw2/errs"
	"hw2/repository"
	"strings"
)

// AuthService отвечает за регистрацию и аутентификацию пользователей
type AuthService struct {
	userRepo repository.UserRepository
}

// NewAuthService создает новый сервис аутентификации
func NewAuthService(userRepo repository.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

// Register регистрирует нового пользователя и возвращает JWT токен
func (a *AuthService) Register(username string, password string) (string, error) {
	username = strings.TrimSpace(username)
	if username == "" || password == "" {
		return "", errs.ErrInvalidInput
	}

	_, err := a.userRepo.GetByUsername(username)
	switch {
	case err == nil:
		return "", errs.ErrUserAlreadyExists
	case errors.Is(err, errs.ErrUserNotFound):
		// ок, пользователя нет, можно регистрировать
	default:
		return "", errs.ErrInternal
	}

	passwordHash, err := auth.HashPassword(password)
	if err != nil {
		return "", errs.ErrInternal
	}

	user := &domain.User{Username: username, PasswordHash: passwordHash}
	err = a.userRepo.Create(user)
	if err != nil {
		return "", errs.ErrInternal
	}

	token, err := auth.GenerateToken(username)
	if err != nil {
		return "", errs.ErrInternal
	}
	return token, nil
}

// Login проверяет учетные данные и возвращает JWT токен
func (a *AuthService) Login(username string, password string) (string, error) {
	username = strings.TrimSpace(username)
	if username == "" || password == "" {
		return "", errs.ErrInvalidInput
	}
	user, err := a.userRepo.GetByUsername(username)
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			return "", errs.ErrInvalidCredentials
		}
		return "", errs.ErrInternal
	}

	err = auth.CheckPassword(password, user.PasswordHash)
	if err != nil {
		return "", errs.ErrInvalidCredentials
	}

	token, err := auth.GenerateToken(username)
	if err != nil {
		return "", errs.ErrInternal
	}
	return token, nil
}
