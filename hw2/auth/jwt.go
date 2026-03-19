// Package auth предоставляет функции для аутентификации:
// - генерация и проверка JWT токенов (jwt.go)
// - хеширование и проверка паролей (password.go)
package auth

import (
	"hw2/errs"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// jwtSecret - секретный ключ для подписи JWT токенов.
// Это костыль, в проде должен храниться в безопасном месте.
var jwtSecret = []byte("super-secret")

// GenerateToken создает JWT токен для пользователя.
// Токен действителен 24 часа.
func GenerateToken(username string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ParseToken проверяет JWT токен и возвращает имя пользователя.
// Возвращает errs.ErrInvalidToken если токен невалидный.
func ParseToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errs.ErrInvalidToken
		}
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return "", errs.ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errs.ErrInvalidToken
	}
	username, ok := claims["username"].(string)
	if !ok || username == "" {
		return "", errs.ErrInvalidToken
	}
	return username, nil

}
