package delivery

import (
	"encoding/json"
	"log"
	"net/http"
	"studybuddy/backend/pkg/auth"
	"studybuddy/backend/pkg/httputil"
	"studybuddy/backend/services/users/usecase"
)

type InterestsHandler struct {
	ListCatalog usecase.ListInterests
	GetMine     usecase.GetMyInterests
	ReplaceMine usecase.ReplaceMyInterests
}

type ReplaceMyInterestsRequest struct {
	InterestIDs []string `json:"interest_ids"`
}

func (h *InterestsHandler) HandleListCatalog(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httputil.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	items, err := h.ListCatalog.ListInterests()
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, "failed to list interests")
		return
	}
	httputil.JSON(w, http.StatusOK, map[string]any{"items": items})
}

func (h *InterestsHandler) HandleGetMyInterests(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httputil.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	userID := auth.UserIDFromContext(r.Context())
	if userID == "" {
		httputil.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	items, err := h.GetMine.GetMyInterests(userID)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, "failed to list interests")
		return
	}
	httputil.JSON(w, http.StatusOK, map[string]any{"items": items})
}

func (h *InterestsHandler) HandleReplaceMyInterests(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		httputil.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	userID := auth.UserIDFromContext(r.Context())
	if userID == "" {
		httputil.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req ReplaceMyInterestsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "failed to parse request")
		return
	}

	items, err := h.ReplaceMine.ReplaceMyInterests(usecase.ReplaceMyInterestsInput{
		UserID:      userID,
		InterestIDs: req.InterestIDs,
	})
	if err != nil {
		log.Printf("HandleReplaceMyInterests: use case error: %v", err)
		if err == usecase.ErrInvalidInterestIDs {
			httputil.Error(w, http.StatusBadRequest, "invalid interest ids")
			return
		}
		httputil.Error(w, http.StatusInternalServerError, "failed to update interests")
		return
	}
	httputil.JSON(w, http.StatusOK, map[string]any{"items": items})
}
