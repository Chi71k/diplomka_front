package delivery

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"studybuddy/backend/pkg/auth"
	"studybuddy/backend/pkg/httputil"
	"studybuddy/backend/services/users/domain"
	"studybuddy/backend/services/users/usecase"
)

// UsersHandler exposes user profile HTTP endpoints.
type UsersHandler struct {
	GetMe           usecase.GetMe
	UpdateMe        usecase.UpdateMe
	DeleteMe        usecase.DeleteMe
	SearchServiceURL string // e.g. "http://localhost:8083" — empty disables indexing
}

// UserProfileResponse matches OpenAPI UserProfile (minimal).
type UserProfileResponse struct {
	ID        string            `json:"id"`
	Email     string            `json:"email"`
	FirstName string            `json:"firstName"`
	LastName  string            `json:"lastName"`
	Bio       string            `json:"bio"`
	AvatarURL string            `json:"avatarUrl"`
	Interests []domain.Interest `json:"interests"`
}

// UpdateProfileRequest matches OpenAPI UpdateProfileRequest.
type UpdateProfileRequest struct {
	FirstName *string `json:"firstName"`
	LastName  *string `json:"lastName"`
	Bio       *string `json:"bio"`
	AvatarURL *string `json:"avatarUrl"`
}

// HandleGetMe GET /api/v1/users/me (requires JWT).
func (h *UsersHandler) HandleGetMe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httputil.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	userID := auth.UserIDFromContext(r.Context())
	if userID == "" {
		httputil.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	profile, err := h.GetMe.GetMe(userID)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, "failed to get profile")
		return
	}
	if profile == nil {
		httputil.Error(w, http.StatusNotFound, "profile not found")
		return
	}
	httputil.JSON(w, http.StatusOK, profileToResponse(profile))
}

// HandleUpdateMe PUT /api/v1/users/me (requires JWT).
func (h *UsersHandler) HandleUpdateMe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		httputil.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	userID := auth.UserIDFromContext(r.Context())
	if userID == "" {
		httputil.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var req UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid body")
		return
	}
	in := usecase.UpdateMeInput{UserID: userID}
	in.FirstName = req.FirstName
	in.LastName = req.LastName
	in.Bio = req.Bio
	in.AvatarURL = req.AvatarURL
	profile, err := h.UpdateMe.UpdateMe(in)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, "failed to update profile")
		return
	}

	// Fire-and-forget: re-index the user in the search service.
	// Never blocks the main response; errors are only logged.
	if h.SearchServiceURL != "" {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			body, _ := json.Marshal(map[string]string{"user_id": userID})
			req, err := http.NewRequestWithContext(ctx, http.MethodPost, h.SearchServiceURL+"/index", bytes.NewReader(body))
			if err != nil {
				log.Printf("search index: new request: %v", err)
				return
			}
			req.Header.Set("Content-Type", "application/json")
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				log.Printf("search index: %v", err)
				return
			}
			resp.Body.Close()
		}()
	}

	httputil.JSON(w, http.StatusOK, profileToResponse(profile))
}

// HandleDeleteMe DELETE /api/v1/users/me (requires JWT).
func (h *UsersHandler) HandleDeleteMe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		httputil.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	userID := auth.UserIDFromContext(r.Context())
	if userID == "" {
		httputil.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	if err := h.DeleteMe.DeleteMe(userID); err != nil {
		httputil.Error(w, http.StatusInternalServerError, "failed to delete account")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// HandleHealth GET /health
func (h *UsersHandler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httputil.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	httputil.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func profileToResponse(p *domain.Profile) UserProfileResponse {
	return UserProfileResponse{
		ID:        p.UserID,
		Email:     p.Email,
		FirstName: p.FirstName,
		LastName:  p.LastName,
		Bio:       p.Bio,
		AvatarURL: p.AvatarURL,
		Interests: nil, // TODO: load when interests endpoint is added
	}
}
