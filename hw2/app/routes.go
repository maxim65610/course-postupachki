// Package app регистрирует все HTTP маршруты приложения
// и применяет необходимые middleware.
package app

import (
	"net/http"

	"hw2/handler"
	"hw2/middleware"
)

// registerRoutes регистрирует все эндпоинты API в переданном ServeMux.
//
// Маршруты разделены на три группы:
//  1. Аутентификация (публичные)
//  2. Криптовалюты (требуют аутентификации)
//  3. Расписание обновления (требует аутентификации)
//
// Параметры:
//   - mux: HTTP мультиплексор для регистрации маршрутов
//   - authHandler: обработчик для эндпоинтов аутентификации
//   - cryptoHandler: обработчик для эндпоинтов криптовалют
//   - scheduleHandler: обработчик для эндпоинтов расписания
//
// Все маршруты кроме /auth/* защищены middleware.Auth,
// который проверяет JWT токен в заголовке Authorization.
func registerRoutes(
	mux *http.ServeMux,
	authHandler *handler.AuthHandler,
	cryptoHandler *handler.CryptoHandler,
	scheduleHandler *handler.ScheduleHandler,
) {
	mux.HandleFunc("POST /auth/register", authHandler.Register)
	mux.HandleFunc("POST /auth/login", authHandler.Login)

	mux.Handle("GET /crypto", middleware.Auth(http.HandlerFunc(cryptoHandler.GetAll)))
	mux.Handle("POST /crypto", middleware.Auth(http.HandlerFunc(cryptoHandler.Add)))
	mux.Handle("GET /crypto/{symbol}", middleware.Auth(http.HandlerFunc(cryptoHandler.GetOne)))
	mux.Handle("PUT /crypto/{symbol}/refresh", middleware.Auth(http.HandlerFunc(cryptoHandler.Refresh)))
	mux.Handle("GET /crypto/{symbol}/history", middleware.Auth(http.HandlerFunc(cryptoHandler.GetHistory)))
	mux.Handle("GET /crypto/{symbol}/stats", middleware.Auth(http.HandlerFunc(cryptoHandler.GetStats)))
	mux.Handle("DELETE /crypto/{symbol}", middleware.Auth(http.HandlerFunc(cryptoHandler.Delete)))

	mux.Handle("GET /schedule", middleware.Auth(http.HandlerFunc(scheduleHandler.Get)))
	mux.Handle("PUT /schedule", middleware.Auth(http.HandlerFunc(scheduleHandler.Update)))
	mux.Handle("POST /schedule/trigger", middleware.Auth(http.HandlerFunc(scheduleHandler.Trigger)))
}
