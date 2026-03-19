package domain

// User представляет пользователя системы.
type User struct {
	Username     string // Имя пользователя (уникальное)
	PasswordHash string // Хеш пароля (bcrypt)
}
