package handler

import (
	"hw2/service"
	"net/http"
	"time"
)

// ScheduleHandler обрабатывает запросы для управления расписанием обновления цен.
type ScheduleHandler struct {
	scheduleService *service.ScheduleService
}

// NewScheduleHandler создает новый экземпляр ScheduleHandler.
func NewScheduleHandler(scheduleService *service.ScheduleService) *ScheduleHandler {
	return &ScheduleHandler{scheduleService: scheduleService}
}

type updateScheduleRequest struct {
	Enabled         bool `json:"enabled"`
	IntervalSeconds int  `json:"interval_seconds"`
}
type scheduleResponse struct {
	Enabled         bool   `json:"enabled"`
	IntervalSeconds int    `json:"interval_seconds"`
	LastUpdate      string `json:"last_update"`
	NextUpdate      string `json:"next_update"`
}
type updateScheduleResponse struct {
	Enabled         bool `json:"enabled"`
	IntervalSeconds int  `json:"interval_seconds"`
}

type triggerScheduleResponse struct {
	UpdatedCount int    `json:"updated_count"`
	Timestamp    string `json:"timestamp"`
}

// Get обрабатывает GET /schedule.
func (h *ScheduleHandler) Get(w http.ResponseWriter, r *http.Request) {
	settings, err := h.scheduleService.GetSettings()
	if err != nil {
		writeHandlerError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, scheduleResponse{
		Enabled:         settings.Enabled,
		IntervalSeconds: settings.IntervalSeconds,
		LastUpdate:      settings.LastUpdate.UTC().Format(time.RFC3339),
		NextUpdate:      settings.NextUpdate.UTC().Format(time.RFC3339),
	})
}

// Update обрабатывает PUT /schedule.
func (handler *ScheduleHandler) Update(w http.ResponseWriter, r *http.Request) {
	var req updateScheduleRequest
	if err := decodeJSON(r, &req); err != nil {
		writeHandlerError(w, err)
		return
	}

	settings, err := handler.scheduleService.UpdateSettings(req.Enabled, req.IntervalSeconds)
	if err != nil {
		writeHandlerError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, updateScheduleResponse{
		Enabled:         settings.Enabled,
		IntervalSeconds: settings.IntervalSeconds,
	})
}

// Trigger обрабатывает POST /schedule/trigger.
func (handler *ScheduleHandler) Trigger(w http.ResponseWriter, r *http.Request) {
	result, err := handler.scheduleService.TriggerNow()
	if err != nil {
		writeHandlerError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, triggerScheduleResponse{
		UpdatedCount: result.UpdatedCount,
		Timestamp:    result.Timestamp.UTC().Format(time.RFC3339),
	})
}
