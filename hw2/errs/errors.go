// Package errs содержит глобальные ошибки приложения.
// Используется для единообразной обработки ошибок во всех слоях.
package errs

import "errors"

// Ошибки аутентификации и пользователей
var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
)

// Ошибки работы с криптовалютами
var (
	ErrCryptoNotFound      = errors.New("crypto not found")
	ErrCryptoAlreadyExists = errors.New("crypto already exists")
)

// Ошибки расписания и валидации
var (
	ErrScheduleNotFound = errors.New("schedule not found")
	ErrInvalidInterval  = errors.New("invalid interval")
	ErrInvalidInput     = errors.New("invalid input")
)

// Системные ошибки
var (
	ErrExternalService = errors.New("external service error")
	ErrInternal        = errors.New("internal server error")
)
