// Package middleware предоставляет HTTP middleware для аутентификации
package middleware

import (
	"context"
	"hw2/auth"
	"hw2/errs"
	"hw2/handler"
	"net/http"
	"strings"
)

type contextKey string

// ключ для username в контексте
const UsernameContextKey contextKey = "username"

// Auth проверяет JWT токен из заголовка Authorization
// При успехе добавляет username в контекст
// При ошибке возвращает 401
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := strings.TrimSpace(r.Header.Get("Authorization"))
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			handler.WriteError(w, http.StatusUnauthorized, errs.ErrInvalidToken)
			return
		}

		token := strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))
		if token == "" {
			handler.WriteError(w, http.StatusUnauthorized, errs.ErrInvalidToken)
			return
		}

		username, err := auth.ParseToken(token)
		if err != nil {
			handler.WriteError(w, http.StatusUnauthorized, errs.ErrInvalidToken)
			return
		}
		ctx := context.WithValue(r.Context(), UsernameContextKey, username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
