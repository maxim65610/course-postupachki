package repository

import "hw2/domain"

// UserRepository определяет методы для работы с пользователями
type UserRepository interface {
	Create(user *domain.User) error
	GetByUsername(username string) (*domain.User, error)
}
