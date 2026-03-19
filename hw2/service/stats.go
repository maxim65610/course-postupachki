package service

// CryptoStats содержит статистику по ценам криптовалюты
type CryptoStats struct {
	MinPrice           float64 `json:"min_price"`
	MaxPrice           float64 `json:"max_price"`
	AvgPrice           float64 `json:"avg_price"`
	PriceChange        float64 `json:"price_change"`
	PriceChangePercent float64 `json:"price_change_percent"`
	RecordsCount       int     `json:"records_count"`
}

// CryptoStatsResult результат статистики для конкретной криптовалюты
type CryptoStatsResult struct {
	Symbol       string      `json:"symbol"`
	CurrentPrice float64     `json:"current_price"`
	Stats        CryptoStats `json:"stats"`
}
