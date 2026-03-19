package ram_storage

import (
	"hw2/domain"
	"sync"
	"time"
)

// ScheduleRepository хранит настройки расписания в памяти
type ScheduleRepository struct {
	settings domain.ScheduleSettings
	mu       sync.RWMutex
}

// NewScheduleRepository создает репозиторий с настройками по умолчанию
func NewScheduleRepository() *ScheduleRepository {
	return &ScheduleRepository{
		settings: domain.ScheduleSettings{
			Enabled:         true,
			IntervalSeconds: 30,
			LastUpdate:      time.Time{},
			NextUpdate:      time.Now().UTC().Add(30 * time.Second),
		},
	}
}

// Get возвращает текущие настройки
func (r *ScheduleRepository) Get() (domain.ScheduleSettings, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.settings, nil
}

// Update обновляет настройки
func (r *ScheduleRepository) Update(settings domain.ScheduleSettings) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.settings = settings
	return nil
}
