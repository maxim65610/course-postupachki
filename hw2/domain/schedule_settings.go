package domain

import "time"

// ScheduleSettings представляет настройки автообновления цен.
type ScheduleSettings struct {
	Enabled         bool      // Включено ли автообновление
	IntervalSeconds int       // Интервал обновления в секундах
	LastUpdate      time.Time // Время последнего обновления
	NextUpdate      time.Time // Время следующего обновления
}
