// Package domain содержит модели данных предметной области.
// Все структуры представляют основные сущности приложения.
package domain

import "time"

// Crypto представляет криптовалюту для отслеживания.
type Crypto struct {
	Symbol       string        // Тикер
	Name         string        // Полное название
	CoinGeckoID  string        // ID в CoinGecko API
	CurrentPrice float64       // Текущая цена в USD
	LastUpdated  time.Time     // Время последнего обновления
	History      []PriceRecord // История цен
}
