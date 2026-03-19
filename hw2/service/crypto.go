package service

import (
	"errors"
	"hw2/coingecko"
	"hw2/domain"
	"hw2/errs"
	"hw2/repository"
	"strings"
	"time"
)

// CryptoService отвечает за операции с криптовалютами
type CryptoService struct {
	repo   repository.CryptoRepository
	client coingecko.PriceProvider
}

// NewCryptoService создает новый сервис для работы с криптовалютами
func NewCryptoService(repo repository.CryptoRepository, priceProvider coingecko.PriceProvider) *CryptoService {
	return &CryptoService{repo: repo, client: priceProvider}

}

// GetAllCryptos возвращает все отслеживаемые криптовалюты
func (s *CryptoService) GetAllCryptos() ([]*domain.Crypto, error) {
	cryptos, err := s.repo.GetAll()
	if err != nil {
		return nil, errs.ErrInternal
	}
	return cryptos, nil
}

// GetCrypto возвращает криптовалюту по символу
func (s *CryptoService) GetCrypto(symbol string) (*domain.Crypto, error) {
	symbol = normalizeSymbol(symbol)
	if symbol == "" {
		return nil, errs.ErrInvalidInput
	}
	return s.repo.Get(symbol)
}

// DeleteCrypto удаляет криптовалюту из отслеживания
func (s *CryptoService) DeleteCrypto(symbol string) error {

	symbol = normalizeSymbol(symbol)

	if symbol == "" {
		return errs.ErrInvalidInput
	}

	err := s.repo.Delete(symbol)

	return err
}

// GetHistory возвращает историю цен криптовалюты
func (s *CryptoService) GetHistory(symbol string) ([]domain.PriceRecord, error) {
	symbol = normalizeSymbol(symbol)
	if symbol == "" {
		return nil, errs.ErrInvalidInput
	}

	crypto, err := s.repo.Get(symbol)
	if err != nil {
		return nil, errs.ErrInternal
	}

	return crypto.History, nil
}

// AddCrypto добавляет новую криптовалюту для отслеживания
func (s *CryptoService) AddCrypto(symbol string) (*domain.Crypto, error) {
	symbol = normalizeSymbol(symbol)

	if symbol == "" {
		return nil, errs.ErrInvalidInput
	}

	_, err := s.repo.Get(symbol)
	if err == nil {
		return nil, errs.ErrCryptoAlreadyExists
	}
	if !errors.Is(err, errs.ErrCryptoNotFound) {
		return nil, errs.ErrInternal
	}
	coinID, coinName, err := s.client.ResolveSymbol(symbol)
	if err != nil {
		return nil, errs.ErrExternalService
	}
	price := 0.0
	if p, err := s.client.GetCurrentPrice(coinID); err == nil {
		price = p
	}

	now := time.Now().UTC()

	crypto := &domain.Crypto{
		Symbol:       symbol,
		Name:         coinName,
		CoinGeckoID:  coinID,
		CurrentPrice: price,
		LastUpdated:  now,
		History: []domain.PriceRecord{
			{
				Price:     price,
				Timestamp: now,
			},
		},
	}
	err = s.repo.Add(crypto)
	if err != nil {
		return nil, errs.ErrInternal
	}
	return crypto, nil
}

// RefreshCrypto принудительно обновляет цену криптовалюты
func (s *CryptoService) RefreshCrypto(symbol string) (*domain.Crypto, error) {
	symbol = normalizeSymbol(symbol)
	if symbol == "" {
		return nil, errs.ErrInvalidInput
	}

	crypto, err := s.repo.Get(symbol)
	if err != nil {
		return nil, errs.ErrInternal
	}

	price, err := s.client.GetCurrentPrice(crypto.CoinGeckoID)
	if err != nil {
		return nil, errs.ErrExternalService
	}
	now := time.Now().UTC()

	crypto.CurrentPrice = price
	crypto.LastUpdated = now
	crypto.History = append(crypto.History, domain.PriceRecord{
		Price:     price,
		Timestamp: now,
	})

	if len(crypto.History) > 100 {
		crypto.History = crypto.History[len(crypto.History)-100:]
	}

	if err := s.repo.Update(crypto); err != nil {
		return nil, errs.ErrInternal
	}

	return crypto, nil
}

// GetStats возвращает статистику по ценам криптовалюты
func (s *CryptoService) GetStats(symbol string) (*CryptoStatsResult, error) {
	symbol = normalizeSymbol(symbol)
	if symbol == "" {
		return nil, errs.ErrInvalidInput
	}

	crypto, err := s.repo.Get(symbol)
	if err != nil {
		return nil, errs.ErrInternal
	}

	history := crypto.History
	if len(history) == 0 {
		return &CryptoStatsResult{
			Symbol:       crypto.Symbol,
			CurrentPrice: crypto.CurrentPrice,
			Stats: CryptoStats{
				MinPrice:           0,
				MaxPrice:           0,
				AvgPrice:           0,
				PriceChange:        0,
				PriceChangePercent: 0,
				RecordsCount:       0,
			},
		}, nil
	}

	minPrice := history[0].Price
	maxPrice := history[0].Price
	sum := 0.0

	for _, record := range history {
		if record.Price < minPrice {
			minPrice = record.Price
		}
		if record.Price > maxPrice {
			maxPrice = record.Price
		}
		sum += record.Price
	}

	avgPrice := sum / float64(len(history))

	firstPrice := history[0].Price
	lastPrice := history[len(history)-1].Price
	priceChange := lastPrice - firstPrice

	priceChangePercent := 0.0
	if firstPrice != 0 {
		priceChangePercent = (priceChange / firstPrice) * 100
	}

	return &CryptoStatsResult{
		Symbol:       crypto.Symbol,
		CurrentPrice: crypto.CurrentPrice,
		Stats: CryptoStats{
			MinPrice:           minPrice,
			MaxPrice:           maxPrice,
			AvgPrice:           avgPrice,
			PriceChange:        priceChange,
			PriceChangePercent: priceChangePercent,
			RecordsCount:       len(history),
		},
	}, nil
}

// normalizeSymbol приводит символ к верхнему регистру и убирает пробелы
func normalizeSymbol(symbol string) string {
	return strings.ToUpper(strings.TrimSpace(symbol))

}
