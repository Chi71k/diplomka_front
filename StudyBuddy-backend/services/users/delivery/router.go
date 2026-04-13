package delivery

import (
	"net/http"

	"studybuddy/backend/pkg/auth"
)

// NewRouter returns the users service HTTP router.
// JWT secret must match the Auth service secret.
func NewRouter(h *UsersHandler, ih *InterestsHandler, jwtSecret []byte) http.Handler {
	protect := auth.Middleware(jwtSecret)
	mux := http.NewServeMux()
	// Use path-only patterns for Go < 1.22; method is checked inside handlers.
	mux.HandleFunc("/health", h.HandleHealth)
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
	mux.Handle("/api/v1/interests", protect(http.HandlerFunc(ih.HandleListCatalog)))

	mux.Handle("/api/v1/users/me/interests", protect(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			ih.HandleGetMyInterests(w, r)
		case http.MethodPut:
			ih.HandleReplaceMyInterests(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})))
	return mux
}
