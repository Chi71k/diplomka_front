package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Handler holds the DB pool used by all HTTP handlers.
type Handler struct {
	pool *pgxpool.Pool
}

// HandleIndex POST /index
// Body: {"user_id": "<uuid>"}
// Builds profile text, calls Gemini for an embedding, upserts into user_embeddings.
func (h *Handler) HandleIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req struct {
		UserID string `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if req.UserID == "" {
		writeError(w, http.StatusBadRequest, "user_id required")
		return
	}

	text, err := GetUserProfile(r.Context(), h.pool, req.UserID)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	if text == "" {
		writeError(w, http.StatusUnprocessableEntity, "user profile is empty — nothing to index")
		return
	}

	embedding, err := GetEmbedding(r.Context(), text)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get embedding: "+err.Error())
		return
	}

	if err := UpsertEmbedding(r.Context(), h.pool, req.UserID, embedding); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to upsert embedding: "+err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"status":  "indexed",
		"user_id": req.UserID,
	})
}

// HandleSearch GET /search?q=...&limit=5
// Embeds the query, runs cosine-distance search, returns []UserResult.
func (h *Handler) HandleSearch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	q := r.URL.Query().Get("q")
	if q == "" {
		writeError(w, http.StatusBadRequest, "q required")
		return
	}

	limit := 5
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			if n > 20 {
				n = 20
			}
			limit = n
		}
	}

	embedding, err := GetEmbedding(r.Context(), q)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get embedding: "+err.Error())
		return
	}

	results, err := SearchUsers(r.Context(), h.pool, embedding, limit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "search failed: "+err.Error())
		return
	}

	writeJSON(w, http.StatusOK, results)
}

// HandleHealth GET /health
func (h *Handler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// writeJSON writes a JSON response with the given status code.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if v != nil {
		_ = json.NewEncoder(w).Encode(v)
	}
}

// writeError writes a JSON error response.
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
