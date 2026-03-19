package domain

import "time"

// PriceRecord представляет запись цены в определенный момент времени.
type PriceRecord struct {
	Price     float64   // Цена в USD
	Timestamp time.Time // Время записи
}
