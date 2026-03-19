package auth

import "golang.org/x/crypto/bcrypt"

// HashPassword Хеширует пароль с помощью bcrypt.
// Возвращает хеш или ошибку при неудаче.
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// CheckPassword Проверяет соответствие пароля и хеша.
func CheckPassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
