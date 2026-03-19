// Package coingecko предоставляет клиент для работы с CoinGecko API.
// Реализует получение списка монет и текущих цен криптовалют.
package coingecko

import (
	"encoding/json"
	"hw2/errs"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const baseURL = "https://api.coingecko.com/api/v3"

// Coin представляет криптовалюту в ответе CoinGecko API.
type Coin struct {
	ID     string `json:"id"`
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
}

// PriceProvider определяет интерфейс для получения данных о криптовалютах.
type PriceProvider interface {
	ResolveSymbol(symbol string) (string, string, error)
	GetCurrentPrice(id string) (float64, error)
}

// Client представляет HTTP клиент для CoinGecko API.
type Client struct {
	symbolMap  map[string]Coin // Кэш маппинга symbol -> Coin
	httpClient *http.Client    // HTTP клиент с таймаутом
}

// NewClient создает новый клиент и загружает список всех монет.
// Возвращает ошибку errs.ErrExternalService при проблемах с API.
func NewClient() (*Client, error) {
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := httpClient.Get(baseURL + "/coins/list")
	if err != nil {
		return nil, errs.ErrExternalService
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errs.ErrExternalService
	}

	var coins []Coin
	if err := json.NewDecoder(resp.Body).Decode(&coins); err != nil {
		return nil, errs.ErrExternalService
	}

	symbolMap := make(map[string]Coin)
	for _, coin := range coins {
		symbol := strings.ToUpper(strings.TrimSpace(coin.Symbol))
		if symbol == "" || coin.ID == "" {
			continue
		}

		if _, exists := symbolMap[symbol]; !exists {
			symbolMap[symbol] = coin
		}
	}

	return &Client{
		symbolMap:  symbolMap,
		httpClient: httpClient,
	}, nil
}

// ResolveSymbol преобразует тикер в CoinGecko ID и полное название.
// Возвращает errs.ErrInvalidInput для пустого символа,
// errs.ErrCryptoNotFound если символ не найден.
func (c *Client) ResolveSymbol(symbol string) (string, string, error) {
	symbol = strings.ToUpper(strings.TrimSpace(symbol))
	if symbol == "" {
		return "", "", errs.ErrInvalidInput
	}

	coin, ok := c.symbolMap[symbol]
	if !ok {
		return "", "", errs.ErrCryptoNotFound
	}

	return coin.ID, coin.Name, nil
}

// GetCurrentPrice получает текущую цену в USD по CoinGecko ID.
// Возвращает errs.ErrInvalidInput для пустого ID,
// errs.ErrCryptoNotFound если монета не найдена,
// errs.ErrExternalService при ошибках API.
func (c *Client) GetCurrentPrice(id string) (float64, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return 0, errs.ErrInvalidInput
	}

	u, err := url.Parse(baseURL + "/simple/price")
	if err != nil {
		return 0, errs.ErrInternal
	}

	q := u.Query()
	q.Set("ids", id)
	q.Set("vs_currencies", "usd")
	u.RawQuery = q.Encode()

	resp, err := c.httpClient.Get(u.String())
	if err != nil {
		return 0, errs.ErrExternalService
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return 0, errs.ErrCryptoNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return 0, errs.ErrExternalService
	}

	var data map[string]map[string]float64
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, errs.ErrExternalService
	}

	coinData, ok := data[id]
	if !ok {
		return 0, errs.ErrCryptoNotFound
	}

	price, ok := coinData["usd"]
	if !ok {
		return 0, errs.ErrCryptoNotFound
	}

	return price, nil
}
