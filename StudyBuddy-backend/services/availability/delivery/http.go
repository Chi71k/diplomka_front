package delivery

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"studybuddy/backend/pkg/auth"
	"studybuddy/backend/pkg/httputil"
	"studybuddy/backend/services/availability/domain"
	"studybuddy/backend/services/availability/usecase"
	"time"
)

// AvailabilityHandler exposes all availability HTTP endpoints.
type AvailabilityHandler struct {
	ListSlots      usecase.ListSlots
	CreateSlot     usecase.CreateSlot
	DeleteSlot     usecase.DeleteSlot
	GCalConnect    usecase.GCalConnect
	GCalImport     usecase.GCalImport
	GCalDisconnect usecase.GCalDisconnect
}

// request / response shapes

// SlotResponse is the API shape returned for every slot.
type SlotResponse struct {
	ID        string `json:"id"`
	DayOfWeek int    `json:"dayOfWeek"` // 0 = Monday … 6 = Sunday
	StartTime string `json:"startTime"` // "HH:MM"
	EndTime   string `json:"endTime"`   // "HH:MM"
	Timezone  string `json:"timezone"`
}

// CreateSlotRequest is the JSON body for POST /api/v1/availability/slots.
type CreateSlotRequest struct {
	DayOfWeek int    `json:"dayOfWeek"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	Timezone  string `json:"timezone"`
}

// DisconnectRequest is the optional JSON body for DELETE /api/v1/availability/gcal/disconnect.
type DisconnectRequest struct {
	// DeleteSlots removes all imported slots when true.
	// Defaults to false so the user keeps their schedule after disconnecting.
	DeleteSlots bool `json:"deleteSlots"`
}

// health

// HandleHealth GET /health
func (h *AvailabilityHandler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httputil.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	httputil.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// ── slots collection: GET + POST /api/v1/availability/slots ──────────────────

// HandleSlotsCollection dispatches GET and POST on /api/v1/availability/slots.
func (h *AvailabilityHandler) HandleSlotsCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleListSlots(w, r)
	case http.MethodPost:
		h.handleCreateSlot(w, r)
	default:
		httputil.Error(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

// handleListSlots GET /api/v1/availability/slots
func (h *AvailabilityHandler) handleListSlots(w http.ResponseWriter, r *http.Request) {
	userID := auth.UserIDFromContext(r.Context())
	if userID == "" {
		httputil.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	slots, err := h.ListSlots.ListSlots(userID)
	if err != nil {
		log.Printf("handleListSlots: %v", err)
		httputil.Error(w, http.StatusInternalServerError, "failed to list slots")
		return
	}

	resp := make([]SlotResponse, 0, len(slots))
	for _, s := range slots {
		resp = append(resp, slotToResponse(s))
	}
	httputil.JSON(w, http.StatusOK, map[string]any{"items": resp})
}

// handleCreateSlot POST /api/v1/availability/slots
func (h *AvailabilityHandler) handleCreateSlot(w http.ResponseWriter, r *http.Request) {
	userID := auth.UserIDFromContext(r.Context())
	if userID == "" {
		httputil.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req CreateSlotRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.StartTime == "" || req.EndTime == "" || req.Timezone == "" {
		httputil.Error(w, http.StatusBadRequest, "startTime, endTime, timezone are required")
		return
	}

	out, err := h.CreateSlot.CreateSlot(usecase.CreateSlotInput{
		UserID:    userID,
		DayOfWeek: req.DayOfWeek,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Timezone:  req.Timezone,
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidDayOfWeek),
			errors.Is(err, domain.ErrInvalidTimezone),
			errors.Is(err, domain.ErrInvalidTimeFormat),
			errors.Is(err, domain.ErrInvalidTimeRange):
			httputil.Error(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, domain.ErrSlotConflict):
			httputil.Error(w, http.StatusConflict, "slot overlaps with an existing slot")
		default:
			log.Printf("handleCreateSlot: %v", err)
			httputil.Error(w, http.StatusInternalServerError, "failed to create slot")
		}
		return
	}

	httputil.JSON(w, http.StatusCreated, slotToResponse(out.Slot))
}

// ── slot item: DELETE /api/v1/availability/slots/{id} ────────────────────────

// HandleSlotItem dispatches DELETE on /api/v1/availability/slots/{id}.
func (h *AvailabilityHandler) HandleSlotItem(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/availability/slots/")
	if id == "" || strings.Contains(id, "/") {
		httputil.Error(w, http.StatusNotFound, "not found")
		return
	}

	switch r.Method {
	case http.MethodDelete:
		h.handleDeleteSlot(w, r, id)
	default:
		httputil.Error(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

// handleDeleteSlot DELETE /api/v1/availability/slots/{id}
func (h *AvailabilityHandler) handleDeleteSlot(w http.ResponseWriter, r *http.Request, slotID string) {
	userID := auth.UserIDFromContext(r.Context())
	if userID == "" {
		httputil.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	err := h.DeleteSlot.DeleteSlot(usecase.DeleteSlotInput{
		UserID: userID,
		SlotID: slotID,
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrSlotNotFound):
			httputil.Error(w, http.StatusNotFound, "slot not found")
		case errors.Is(err, domain.ErrForbidden):
			httputil.Error(w, http.StatusForbidden, "you do not own this slot")
		default:
			log.Printf("handleDeleteSlot: %v", err)
			httputil.Error(w, http.StatusInternalServerError, "failed to delete slot")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ── gcal connect: GET /api/v1/availability/gcal/connect ──────────────────────

// HandleGCalConnect GET /api/v1/availability/gcal/connect
// Returns the Google OAuth consent URL. The frontend redirects the user there.
func (h *AvailabilityHandler) HandleGCalConnect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httputil.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	userID := auth.UserIDFromContext(r.Context())
	if userID == "" {
		httputil.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	authURL, err := h.GCalConnect.GetAuthUrl(userID)
	if err != nil {
		log.Printf("HandleGCalConnect: build auth url: %v", err)
		httputil.Error(w, http.StatusInternalServerError, "failed to build google auth url")
		return
	}

	httputil.JSON(w, http.StatusOK, map[string]string{"authUrl": authURL})
}

// ── gcal callback: GET /api/v1/availability/gcal/callback ────────────────────

// HandleGCalCallback GET /api/v1/availability/gcal/callback
// Google redirects the user's browser here after they approve or deny consent.
// This endpoint is NOT protected by JWT middleware — identity is verified via
// the HMAC-signed state parameter that the GCalConnect use case produces.
func (h *AvailabilityHandler) HandleGCalCallback(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httputil.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	q := r.URL.Query()

	// Google sends an "error" query param when the user denies consent.
	if errParam := q.Get("error"); errParam != "" {
		log.Printf("HandleGCalCallback: google returned error: %s", errParam)
		httputil.Error(w, http.StatusBadRequest, "google calendar access was denied")
		return
	}

	code := q.Get("code")
	state := q.Get("state")

	if code == "" || state == "" {
		httputil.Error(w, http.StatusBadRequest, "missing code or state parameter")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := h.GCalConnect.HandleCallback(ctx, state, code); err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidOAuthState):
			httputil.Error(w, http.StatusBadRequest, "invalid oauth state — please start the connection flow again")
		case errors.Is(err, domain.ErrOAuthStateExpired):
			httputil.Error(w, http.StatusBadRequest, "oauth state expired — please start the connection flow again")
		default:
			log.Printf("HandleGCalCallback: %v", err)
			httputil.Error(w, http.StatusInternalServerError, "failed to complete google calendar connection")
		}
		return
	}

	// The callback is a browser redirect, so respond with a plain success page
	// rather than JSON. The frontend can poll /gcal/status or listen for a
	// postMessage from this window to know when the flow is complete.
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`<!DOCTYPE html>
<html>
<head><title>Connected</title></head>
<body>
  <p>Google Calendar connected successfully. You can close this window.</p>
  <script>
    // Notify the opener tab (if this was opened as a popup) then close.
    if (window.opener) {
      window.opener.postMessage({ type: "GCAL_CONNECTED" }, "*");
      window.close();
    }
  </script>
</body>
</html>`))
}

// ── gcal import: POST /api/v1/availability/gcal/import ───────────────────────

// HandleGCalImport POST /api/v1/availability/gcal/import
// Triggers a manual sync: fetches the user's Google Calendar events for the
// next 4 weeks and converts them into recurring availability slots.
func (h *AvailabilityHandler) HandleGCalImport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httputil.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	userID := auth.UserIDFromContext(r.Context())
	if userID == "" {
		httputil.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	slots, err := h.GCalImport.ImportFromGCal(r.Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrGCalNotConnected):
			httputil.Error(w, http.StatusPreconditionRequired, "google calendar is not connected — call /gcal/connect first")
		case errors.Is(err, domain.ErrGCalSyncDisabled):
			httputil.Error(w, http.StatusForbidden, "google calendar sync is disabled for your account")
		default:
			log.Printf("HandleGCalImport: %v", err)
			httputil.Error(w, http.StatusInternalServerError, "failed to import from google calendar")
		}
		return
	}

	resp := make([]SlotResponse, 0, len(slots))
	for _, s := range slots {
		resp = append(resp, slotToResponse(s))
	}
	httputil.JSON(w, http.StatusOK, map[string]any{"items": resp, "imported": len(resp)})
}

// ── gcal disconnect: DELETE /api/v1/availability/gcal/disconnect ──────────────

// HandleGCalDisconnect DELETE /api/v1/availability/gcal/disconnect
// Removes the stored OAuth connection. An optional JSON body controls whether
// imported slots are also deleted (defaults to false).
func (h *AvailabilityHandler) HandleGCalDisconnect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		httputil.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	userID := auth.UserIDFromContext(r.Context())
	if userID == "" {
		httputil.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Body is optional — if absent or malformed we use the safe default (keep slots).
	var req DisconnectRequest
	_ = json.NewDecoder(r.Body).Decode(&req) // intentionally ignoring decode error

	if err := h.GCalDisconnect.Disconnect(usecase.GCalDisconnectInput{
		UserID:              userID,
		DeleteImportedSlots: req.DeleteSlots,
	}); err != nil {
		log.Printf("HandleGCalDisconnect: %v", err)
		httputil.Error(w, http.StatusInternalServerError, "failed to disconnect google calendar")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// helpers

// slotToResponse converts a domain.Slot into the API response shape.
// StartTime and EndTime are formatted as "HH:MM" in the slot's own timezone.
func slotToResponse(s domain.Slot) SlotResponse {
	return SlotResponse{
		ID:        s.ID,
		DayOfWeek: s.DayOfWeek,
		StartTime: s.StartTime.Format("15:04"),
		EndTime:   s.EndTime.Format("15:04"),
		Timezone:  s.Timezone,
	}
}
