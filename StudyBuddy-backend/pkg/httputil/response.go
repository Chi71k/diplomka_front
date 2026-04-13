package httputil

import (
	"encoding/json"
	"net/http"
)

// JSON writes a JSON response with status code.
func JSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if v != nil {
		_ = json.NewEncoder(w).Encode(v)
	}
}

// ErrBody is a common error response shape.
type ErrBody struct {
	Error string `json:"error"`
}

// Error writes a JSON error response.
func Error(w http.ResponseWriter, status int, message string) {
	JSON(w, status, ErrBody{Error: message})
}
