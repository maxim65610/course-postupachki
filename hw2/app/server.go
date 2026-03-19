// Package app предоставляет основной компонент сервера приложения
// и его конфигурацию.
package app

import (
	"fmt"
	"hw2/coingecko"
	"hw2/handler"
	"hw2/repository/ram_storage"
	"hw2/service"
	"net/http"
)

// Server представляет основной HTTP сервер приложения.
// Содержит в себе http.Server и управляет его жизненным циклом.
type Server struct {
	httpServer *http.Server
}

// NewServer создает и инициализирует новый экземпляр сервера.
//
// Процесс инициализации включает:
//  1. Создание in-memory репозиториев (User, Crypto, Schedule)
//  2. Инициализацию клиента CoinGecko API
//  3. Создание сервисного слоя (Auth, Crypto, Schedule)
//  4. Создание HTTP обработчиков
//  5. Регистрацию всех маршрутов
//  6. Запуск фонового процесса автообновления цен
//
// Возвращаемые значения:
//   - *Server: указатель на созданный сервер
//   - error: ошибка при создании клиента CoinGecko или запуске планировщика
func NewServer() (*Server, error) {
	userRepo := ram_storage.NewUserRepository()
	cryptoRepo := ram_storage.NewCryptoRepository()
	scheduleRepo := ram_storage.NewScheduleRepository()

	coinGeckoClient, err := coingecko.NewClient()
	if err != nil {
		return nil, fmt.Errorf("Error creating coingecko client: %v", err)
	}

	authService := service.NewAuthService(userRepo)
	cryptoService := service.NewCryptoService(cryptoRepo, coinGeckoClient)
	scheduleService := service.NewScheduleService(scheduleRepo, cryptoService)

	authHandler := handler.NewAuthHandler(authService)
	cryptoHandler := handler.NewCryptoHandler(cryptoService)
	scheduleHandler := handler.NewScheduleHandler(scheduleService)

	mux := http.NewServeMux()
	registerRoutes(mux, authHandler, cryptoHandler, scheduleHandler)

	if err := scheduleService.Start(); err != nil {
		return nil, err
	}

	return &Server{
		httpServer: &http.Server{
			Addr:    ":8080",
			Handler: mux,
		},
	}, nil
}

// Run запускает HTTP сервер и начинает прослушивание входящих запросов.
func (s *Server) Run() error {
	fmt.Printf("Listening on port 8080\n")
	return s.httpServer.ListenAndServe()
}
