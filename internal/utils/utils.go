package utils

import (
	"encoding/json"
	"net/http"

	"github.com/viniferr33/img-processor/internal/constants"
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

func SplitBearerToken(authHeader string) (string, bool) {
	const prefix = "Bearer "
	if len(authHeader) > len(prefix) && authHeader[:len(prefix)] == prefix {
		return authHeader[len(prefix):], true
	}
	return "", false
}

func GetUserIDFromContext(r *http.Request) (string, bool) {
	if userID, ok := r.Context().Value(constants.ContextKeyUserID).(string); ok {
		return userID, true
	}
	return "", false
}
