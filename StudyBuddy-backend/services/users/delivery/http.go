package delivery

import (
	"encoding/json"
	"net/http"

	"studybuddy/backend/pkg/auth"
	"studybuddy/backend/pkg/httputil"
	"studybuddy/backend/services/users/domain"
	"studybuddy/backend/services/users/usecase"
)

// UsersHandler exposes user profile HTTP endpoints.
type UsersHandler struct {
	GetMe    usecase.GetMe
	UpdateMe usecase.UpdateMe
	DeleteMe usecase.DeleteMe
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

// HandleGetUserByID GET /api/v1/users/:id (requires JWT).
func (h *UsersHandler) HandleGetUserByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httputil.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	const prefix = "/api/v1/users/"
	id := r.URL.Path[len(prefix):]
	if id == "" {
		httputil.Error(w, http.StatusBadRequest, "missing user id")
		return
	}
	profile, err := h.GetMe.GetMe(id)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, "failed to get profile")
		return
	}
	if profile == nil {
		httputil.Error(w, http.StatusNotFound, "user not found")
		return
	}
	httputil.JSON(w, http.StatusOK, profileToResponse(profile))
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
