// Package repository определяет интерфейсы для работы с хранилищем
package repository

import "hw2/domain"

// CryptoRepository определяет методы для работы с криптовалютами
type CryptoRepository interface {
	Add(crypto *domain.Crypto) error
	Get(symbol string) (*domain.Crypto, error)
	GetAll() ([]*domain.Crypto, error)
	Update(crypto *domain.Crypto) error
	Delete(symbol string) error
}
