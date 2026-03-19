// Package ram_storage предоставляет реализации репозиториев для хранения в памяти
package ram_storage

import (
	"hw2/domain"
	"hw2/errs"
	"sync"
)

// CryptoRepository хранит криптовалюты в памяти
type CryptoRepository struct {
	cryptos map[string]*domain.Crypto
	mu      sync.RWMutex
}

// NewCryptoRepository создает новый репозиторий
func NewCryptoRepository() *CryptoRepository {
	return &CryptoRepository{cryptos: make(map[string]*domain.Crypto)}
}

// Add добавляет новую криптовалюту
func (cr *CryptoRepository) Add(crypto *domain.Crypto) error {
	cr.mu.Lock()
	defer cr.mu.Unlock()
	if _, ok := cr.cryptos[crypto.Symbol]; ok {
		return errs.ErrCryptoAlreadyExists
	}
	cr.cryptos[crypto.Symbol] = crypto
	return nil
}

// Get возвращает криптовалюту по символу
func (cr *CryptoRepository) Get(symbol string) (*domain.Crypto, error) {
	cr.mu.RLock()
	defer cr.mu.RUnlock()
	crypto, ok := cr.cryptos[symbol]
	if !ok {
		return nil, errs.ErrCryptoNotFound
	}
	return crypto, nil
}

// GetAll возвращает все криптовалюты
func (cr *CryptoRepository) GetAll() ([]*domain.Crypto, error) {
	cr.mu.RLock()
	defer cr.mu.RUnlock()
	result := make([]*domain.Crypto, 0, len(cr.cryptos))
	for _, v := range cr.cryptos {
		result = append(result, v)
	}
	return result, nil
}

// Update обновляет существующую криптовалюту
func (cr *CryptoRepository) Update(crypto *domain.Crypto) error {
	cr.mu.Lock()
	defer cr.mu.Unlock()
	_, ok := cr.cryptos[crypto.Symbol]
	if !ok {
		return errs.ErrCryptoNotFound
	}
	cr.cryptos[crypto.Symbol] = crypto
	return nil
}

// Delete удаляет криптовалюту по символу
func (cr *CryptoRepository) Delete(symbol string) error {
	cr.mu.Lock()
	defer cr.mu.Unlock()
	if _, ok := cr.cryptos[symbol]; !ok {
		return errs.ErrCryptoNotFound
	}

	delete(cr.cryptos, symbol)
	return nil
}
