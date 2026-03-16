package delivery

import (
	"encoding/json"
	"net/http"
	"time"

	"studybuddy/backend/pkg/httputil"
	"studybuddy/backend/services/auth/domain"
	"studybuddy/backend/services/auth/usecase"
)

// AuthHandler exposes auth HTTP endpoints.
type AuthHandler struct {
	Register usecase.Register
	Login    usecase.Login
}

// RegisterRequest matches OpenAPI RegisterRequest.
type RegisterRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

// LoginRequest matches OpenAPI LoginRequest.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResponse matches OpenAPI AuthResponse.
type AuthResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int64  `json:"expiresIn"`
}

// HandleRegister POST /api/v1/auth/register
func (h *AuthHandler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httputil.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid body")
		return
	}
	if req.Email == "" || req.Password == "" || req.FirstName == "" || req.LastName == "" {
		httputil.Error(w, http.StatusBadRequest, "email, password, firstName, lastName required")
		return
	}
	out, err := h.Register.Register(usecase.RegisterInput{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	})
	if err != nil {
		switch err {
		case domain.ErrEmailExists:
			httputil.Error(w, http.StatusConflict, "email already registered")
			return
		default:
			httputil.Error(w, http.StatusInternalServerError, "registration failed")
			return
		}
	}
	expiresIn := out.ExpiresAt - time.Now().Unix()
	if expiresIn < 0 {
		expiresIn = 0
	}
	httputil.JSON(w, http.StatusCreated, AuthResponse{
		AccessToken:  out.AccessToken,
		RefreshToken: out.RefreshToken,
		ExpiresIn:    expiresIn,
	})
}

// HandleLogin POST /api/v1/auth/login
func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httputil.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid body")
		return
	}
	if req.Email == "" || req.Password == "" {
		httputil.Error(w, http.StatusBadRequest, "email and password required")
		return
	}
	out, err := h.Login.Login(usecase.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		switch err {
		case domain.ErrInvalidCreds, domain.ErrUserInactive:
			httputil.Error(w, http.StatusUnauthorized, "invalid credentials")
			return
		default:
			httputil.Error(w, http.StatusInternalServerError, "login failed")
			return
		}
	}
	expiresIn := out.ExpiresAt - time.Now().Unix()
	if expiresIn < 0 {
		expiresIn = 0
	}
	httputil.JSON(w, http.StatusOK, AuthResponse{
		AccessToken:  out.AccessToken,
		RefreshToken: out.RefreshToken,
		ExpiresIn:    expiresIn,
	})
}

// HandleHealth GET /health
func (h *AuthHandler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httputil.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	httputil.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
