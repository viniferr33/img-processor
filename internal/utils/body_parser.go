package utils

import (
	"encoding/json"
	"net/http"
)

func ParseJsonBody(r *http.Request, dst interface{}) error {
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(dst); err != nil {
		return err
	}

	return nil
}

func WriteJsonResponse(w http.ResponseWriter, statusCode int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(data); err != nil {
		return err
	}
	return nil
}
