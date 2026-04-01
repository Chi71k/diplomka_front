package delivery

import (
	"net/http"
	"strings"

	"studybuddy/backend/pkg/auth"
)

func NewRouter(h *MatchingHandler, jwtSecret []byte) http.Handler {
	protect := auth.Middleware(jwtSecret)
	mux := http.NewServeMux()

	mux.HandleFunc("/health", h.HandleHealth)

	mux.Handle(
		"/api/v1/matching/candidates",
		protect(http.HandlerFunc(h.HandleListCandidates)),
	)

	mux.Handle(
		"/api/v1/matching/requests",
		protect(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodPost:
				h.HandleSendRequest(w, r)
			case http.MethodGet:
				h.HandleListMatches(w, r)
			default:
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
		})),
	)

	mux.Handle(
		"/api/v1/matching/requests/",
		protect(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path
			if strings.HasSuffix(path, "/respond") {
				h.HandleRespond(w, r)
				return
			}
			h.HandleCancel(w, r)
		})),
	)

	return mux
}
