package delivery

import "net/http"

// NewRouter returns the auth service HTTP router.
func NewRouter(h *AuthHandler) http.Handler {
	mux := http.NewServeMux()
	// Use path-only patterns for Go < 1.22; method is checked inside handlers.
	mux.HandleFunc("/health", h.HandleHealth)
	mux.HandleFunc("/api/v1/auth/register", h.HandleRegister)
	mux.HandleFunc("/api/v1/auth/login", h.HandleLogin)
	return mux
}
