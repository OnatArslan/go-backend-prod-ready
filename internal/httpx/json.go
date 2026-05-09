package httpx

import (
	"encoding/json"
	"errors"
	"net/http"
)

var ErrInvalidResponseType = errors.New("invalid json type")

func Write(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	encoder := json.NewEncoder(w)

	err := encoder.Encode(data)
	if err != nil {
		return ErrInvalidResponseType
	}
	return nil
}
