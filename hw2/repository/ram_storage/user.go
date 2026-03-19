package ram_storage

import (
	"hw2/domain"
	"hw2/errs"
	"sync"
)

// UserRepository хранит пользователей в памяти
type UserRepository struct {
	users map[string]*domain.User
	mu    sync.RWMutex
}

// NewUserRepository создает новый репозиторий пользователей
func NewUserRepository() *UserRepository {
	return &UserRepository{users: make(map[string]*domain.User)}
}

// Create добавляет нового пользователя
func (r *UserRepository) Create(user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.users[user.Username]; exists {
		return errs.ErrUserAlreadyExists
	}

	r.users[user.Username] = user
	return nil
}

// GetByUsername возвращает пользователя по имени
func (r *UserRepository) GetByUsername(username string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if user, ok := r.users[username]; ok {
		return user, nil
	}
	return nil, errs.ErrUserNotFound
}
