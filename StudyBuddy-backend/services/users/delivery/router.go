package delivery

import (
	"net/http"

	"studybuddy/backend/pkg/auth"
)

// NewRouter returns the users service HTTP router.
// JWT secret must match the Auth service secret.
func NewRouter(h *UsersHandler, jwtSecret []byte) http.Handler {
	mux := http.NewServeMux()
	// Use path-only patterns for Go < 1.22; method is checked inside handlers.
	mux.HandleFunc("/health", h.HandleHealth)
	// Protected routes
	protect := auth.Middleware(jwtSecret)
	mux.Handle("/api/v1/users/me", protect(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.HandleGetMe(w, r)
		case http.MethodPut:
			h.HandleUpdateMe(w, r)
		case http.MethodDelete:
			h.HandleDeleteMe(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})))
	return mux
}
