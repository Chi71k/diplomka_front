package delivery

import (
	"net/http"
	"studybuddy/backend/pkg/auth"
)

func NewRouter(h *AvailabilityHandler, jwtSecret []byte) http.Handler {
	protect := auth.Middleware(jwtSecret)
	mux := http.NewServeMux()

	// Health — no authentication required.
	mux.HandleFunc("/health", h.HandleHealth)

	// Slots collection (GET + POST) — JWT required.
	mux.Handle(
		"/api/v1/availability/slots",
		protect(http.HandlerFunc(h.HandleSlotsCollection)),
	)

	// Slot item (DELETE /api/v1/availability/slots/{id}) — JWT required.
	// The trailing slash registers a subtree pattern so the mux forwards all
	// paths under /api/v1/availability/slots/ here. The handler extracts the
	// id by trimming the known prefix from r.URL.Path.
	mux.Handle(
		"/api/v1/availability/slots/",
		protect(http.HandlerFunc(h.HandleSlotItem)),
	)

	// Google Calendar — connect and import are JWT-protected.
	// The callback is intentionally public: Google redirects the user's browser
	// there, so no Authorization header is present. Identity is verified inside
	// the handler via the HMAC-signed state parameter.
	mux.Handle(
		"/api/v1/availability/gcal/connect",
		protect(http.HandlerFunc(h.HandleGCalConnect)),
	)
	mux.HandleFunc(
		"/api/v1/availability/gcal/callback",
		h.HandleGCalCallback,
	)
	mux.Handle(
		"/api/v1/availability/gcal/import",
		protect(http.HandlerFunc(h.HandleGCalImport)),
	)
	mux.Handle(
		"/api/v1/availability/gcal/disconnect",
		protect(http.HandlerFunc(h.HandleGCalDisconnect)),
	)

	return mux
}
