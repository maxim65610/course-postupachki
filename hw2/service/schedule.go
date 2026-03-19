package service

import (
	"hw2/domain"
	"hw2/errs"
	"hw2/repository"
	"sync"
	"time"
)

const (
	minIntervalSeconds = 10   // минимальный интервал обновления
	maxIntervalSeconds = 3600 // максимальный интервал обновления
)

// CryptoRefresh определяет методы для обновления криптовалют
type CryptoRefresh interface {
	GetAllCryptos() ([]*domain.Crypto, error)
	RefreshCrypto(symbol string) (*domain.Crypto, error)
}

// TriggerResult результат принудительного обновления
type TriggerResult struct {
	UpdatedCount int       `json:"updated_count"`
	Timestamp    time.Time `json:"timestamp"`
}

// ScheduleService управляет автоматическим обновлением цен
type ScheduleService struct {
	repo      repository.ScheduleRepository
	refresher CryptoRefresh

	mu     sync.Mutex
	stopCh chan struct{}
}

// NewScheduleService создает новый сервис расписания
func NewScheduleService(repo repository.ScheduleRepository,
	refresher CryptoRefresh) *ScheduleService {
	return &ScheduleService{repo: repo, refresher: refresher}
}

// GetSettings возвращает текущие настройки расписания
func (service *ScheduleService) GetSettings() (domain.ScheduleSettings, error) {
	settings, err := service.repo.Get()
	if err != nil {
		return domain.ScheduleSettings{}, errs.ErrInternal
	}

	return settings, nil
}

// UpdateSettings обновляет настройки расписания
func (service *ScheduleService) UpdateSettings(enabled bool, intervalSeconds int) (domain.ScheduleSettings, error) {
	if intervalSeconds < minIntervalSeconds || intervalSeconds > maxIntervalSeconds {
		return domain.ScheduleSettings{}, errs.ErrInvalidInterval
	}

	settings, err := service.GetSettings()
	if err != nil {
		return domain.ScheduleSettings{}, errs.ErrInternal
	}

	settings.Enabled = enabled
	settings.IntervalSeconds = intervalSeconds
	if enabled {
		settings.NextUpdate = time.Now().UTC().Add(time.Duration(intervalSeconds) * time.Second)
	} else {
		settings.NextUpdate = time.Time{}
	}

	if err := service.repo.Update(settings); err != nil {
		return domain.ScheduleSettings{}, errs.ErrInternal
	}
	service.restartWorker(intervalSeconds, enabled)
	return settings, nil
}

// TriggerNow принудительно запускает обновление всех цен
func (service *ScheduleService) TriggerNow() (TriggerResult, error) {
	cryptos, err := service.refresher.GetAllCryptos()
	if err != nil {
		return TriggerResult{}, errs.ErrInternal
	}

	now := time.Now().UTC()
	updatedCount := 0
	for _, crypto := range cryptos {
		if _, err := service.refresher.RefreshCrypto(crypto.Symbol); err != nil {
			continue
		}
		updatedCount++
	}
	settings, err := service.GetSettings()
	if err != nil {
		return TriggerResult{}, errs.ErrInternal
	}
	settings.LastUpdate = now
	if settings.Enabled {
		settings.NextUpdate = now.Add(time.Duration(settings.IntervalSeconds) * time.Second)
	} else {
		settings.NextUpdate = time.Time{}
	}

	if err := service.repo.Update(settings); err != nil {
		return TriggerResult{}, errs.ErrInternal
	}

	return TriggerResult{UpdatedCount: updatedCount, Timestamp: now}, nil

}

// runWorker запускает фоновый процесс для периодического обновления
func (service *ScheduleService) runWorker(intervalSeconds int, stopCh <-chan struct{}) {
	ticker := time.NewTicker(time.Duration(intervalSeconds) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			_, _ = service.TriggerNow()
		case <-stopCh:
			return
		}
	}
}

// Start запускает фоновое обновление если оно включено
func (service *ScheduleService) Start() error {
	settings, err := service.GetSettings()
	if err != nil {
		return errs.ErrInternal
	}
	if !settings.Enabled {
		return nil
	}
	service.startWorker(settings.IntervalSeconds)
	return nil
}

// startWorker запускает воркер с указанным интервалом
func (service *ScheduleService) startWorker(intervalSeconds int) {
	service.mu.Lock()
	defer service.mu.Unlock()
	if service.stopCh != nil {
		return
	}
	stopCh := make(chan struct{})
	service.stopCh = stopCh

	go service.runWorker(intervalSeconds, stopCh)
}

// Stop останавливает фоновое обновление
func (service *ScheduleService) Stop() {
	service.mu.Lock()
	defer service.mu.Unlock()
	if service.stopCh != nil {
		close(service.stopCh)
		service.stopCh = nil
	}
}

// restartWorker перезапускает воркер с новыми настройками
func (service *ScheduleService) restartWorker(intervalSeconds int, enabled bool) {
	service.mu.Lock()
	defer service.mu.Unlock()

	if service.stopCh != nil {
		close(service.stopCh)
		service.stopCh = nil
	}

	if !enabled {
		return
	}

	stopCh := make(chan struct{})
	service.stopCh = stopCh

	go service.runWorker(intervalSeconds, stopCh)
}
