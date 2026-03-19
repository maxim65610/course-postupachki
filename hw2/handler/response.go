package handler

import (
	"encoding/json"
	"net/http"
)

// writeJSON отправляет JSON ответ с указанным статусом.
func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

// writeError отправляет JSON ответ с ошибкой.
func WriteError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"error": err.Error()})
}
