package handler

import (
	"encoding/json"
	"hw2/errs"
	"net/http"
)

// decodeJSON декодирует тело запроса в указанную структуру.
func decodeJSON(r *http.Request, dst any) error {
	if r.Body == nil {
		return errs.ErrInvalidInput
	}
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(dst); err != nil {
		return errs.ErrInvalidInput
	}
	return nil
}
