package handler

import (
	"errors"
	"net/http"

	"hw2/errs"
)

var errorStatusMap = map[error]int{
	errs.ErrInvalidInput:        http.StatusBadRequest,
	errs.ErrUserAlreadyExists:   http.StatusConflict,
	errs.ErrInvalidCredentials:  http.StatusUnauthorized,
	errs.ErrInvalidToken:        http.StatusUnauthorized,
	errs.ErrUserNotFound:        http.StatusNotFound,
	errs.ErrCryptoNotFound:      http.StatusNotFound,
	errs.ErrCryptoAlreadyExists: http.StatusConflict,
	errs.ErrScheduleNotFound:    http.StatusNotFound,
	errs.ErrInvalidInterval:     http.StatusBadRequest,
	errs.ErrExternalService:     http.StatusInternalServerError,
}

// writeHandlerError преобразует ошибку сервиса в HTTP ответ с соответствующим статусом.
func writeHandlerError(w http.ResponseWriter, err error) {
	for targetErr, status := range errorStatusMap {
		if errors.Is(err, targetErr) {
			WriteError(w, status, err)
			return
		}
	}

	WriteError(w, http.StatusInternalServerError, errors.New("internal server error"))
}
