package repository

import "hw2/domain"

// ScheduleRepository определяет методы для работы с настройками расписания
type ScheduleRepository interface {
	Get() (domain.ScheduleSettings, error)
	Update(settings domain.ScheduleSettings) error
}
