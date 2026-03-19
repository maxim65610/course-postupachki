package handler

import (
	"hw2/domain"
	"hw2/service"
	"net/http"
	"strings"
	"time"
)

// CryptoHandler обрабатывает запросы для работы с криптовалютами.
type CryptoHandler struct {
	cryptoService *service.CryptoService
}

// NewCryptoHandler создает новый экземпляр CryptoHandler.
func NewCryptoHandler(cryptoService *service.CryptoService) *CryptoHandler {
	return &CryptoHandler{cryptoService: cryptoService}
}

type addCryptoRequest struct {
	Symbol string `json:"symbol"`
}

type cryptoResponse struct {
	Symbol       string  `json:"symbol"`
	Name         string  `json:"name"`
	CurrentPrice float64 `json:"current_price"`
	LastUpdated  string  `json:"last_updated"`
}

type cryptoListResponse struct {
	Cryptos []cryptoResponse `json:"cryptos"`
}

type cryptoItemResponse struct {
	Crypto cryptoResponse `json:"crypto"`
}

type priceRecordResponse struct {
	Price     float64 `json:"price"`
	Timestamp string  `json:"timestamp"`
}

type cryptoHistoryResponse struct {
	Symbol  string                `json:"symbol"`
	History []priceRecordResponse `json:"history"`
}

// GetAll обрабатывает GET /crypto.
func (handler *CryptoHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	cryptos, err := handler.cryptoService.GetAllCryptos()
	if err != nil {
		writeHandlerError(w, err)
		return
	}

	response := cryptoListResponse{
		Cryptos: make([]cryptoResponse, 0, len(cryptos)),
	}
	for _, crypto := range cryptos {
		response.Cryptos = append(response.Cryptos, toCryptoResponse(crypto))
	}
	writeJSON(w, http.StatusOK, response)
}

// Add обрабатывает POST /crypto.
func (handler *CryptoHandler) Add(w http.ResponseWriter, r *http.Request) {
	var req addCryptoRequest
	if err := decodeJSON(r, &req); err != nil {
		writeHandlerError(w, err)
		return
	}
	crypto, err := handler.cryptoService.AddCrypto(req.Symbol)
	if err != nil {
		writeHandlerError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, cryptoItemResponse{
		Crypto: toCryptoResponse(crypto),
	})
}

// GetOne обрабатывает GET /crypto/{symbol}.
func (handler *CryptoHandler) GetOne(w http.ResponseWriter, r *http.Request) {
	symbol := r.PathValue("symbol")
	crypto, err := handler.cryptoService.GetCrypto(symbol)
	if err != nil {
		writeHandlerError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toCryptoResponse(crypto))
}

// Refresh обрабатывает PUT /crypto/{symbol}/refresh.
func (handler *CryptoHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	symbol := r.PathValue("symbol")
	crypto, err := handler.cryptoService.RefreshCrypto(symbol)
	if err != nil {
		writeHandlerError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, cryptoItemResponse{Crypto: toCryptoResponse(crypto)})
}

// GetHistory обрабатывает GET /crypto/{symbol}/history.
func (handler *CryptoHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
	symbol := r.PathValue("symbol")

	history, err := handler.cryptoService.GetHistory(symbol)
	if err != nil {
		writeHandlerError(w, err)
		return
	}

	response := cryptoHistoryResponse{
		Symbol:  strings.ToUpper(strings.TrimSpace(symbol)),
		History: make([]priceRecordResponse, 0, len(history)),
	}

	for _, record := range history {
		response.History = append(response.History, priceRecordResponse{
			Price:     record.Price,
			Timestamp: record.Timestamp.UTC().Format(time.RFC3339),
		})
	}

	writeJSON(w, http.StatusOK, response)
}

// GetStats обрабатывает GET /crypto/{symbol}/stats.
func (handler *CryptoHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	symbol := r.PathValue("symbol")
	stats, err := handler.cryptoService.GetStats(symbol)
	if err != nil {
		writeHandlerError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, stats)
}

// Delete обрабатывает DELETE /crypto/{symbol}.
func (handler *CryptoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	symbol := r.PathValue("symbol")
	err := handler.cryptoService.DeleteCrypto(symbol)
	if err != nil {
		writeHandlerError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{})
}

func toCryptoResponse(crypto *domain.Crypto) cryptoResponse {
	return cryptoResponse{
		Symbol:       crypto.Symbol,
		Name:         crypto.Name,
		CurrentPrice: crypto.CurrentPrice,
		LastUpdated:  crypto.LastUpdated.UTC().Format(time.RFC3339),
	}
}
