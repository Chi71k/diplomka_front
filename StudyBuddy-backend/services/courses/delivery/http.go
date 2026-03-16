package delivery

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"studybuddy/backend/pkg/auth"
	"studybuddy/backend/pkg/httputil"
	"studybuddy/backend/services/courses/domain"
	"studybuddy/backend/services/courses/usecase"
)

// CoursesHandler exposes courses HTTP endpoints.
type CoursesHandler struct {
	List   usecase.ListCourses
	Get    usecase.GetCourse
	Create usecase.CreateCourse
	Update usecase.UpdateCourse
	Delete usecase.DeleteCourse
}

// CourseResponse is the API shape for a course.
type CourseResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Subject     string `json:"subject"`
	Level       string `json:"level"`
	OwnerUserID string `json:"ownerUserId"`
}

// CreateCourseRequest is the body for creating a course.
type CreateCourseRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Subject     string `json:"subject"`
	Level       string `json:"level"`
}

// UpdateCourseRequest is the body for partially updating a course.
type UpdateCourseRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Subject     *string `json:"subject"`
	Level       *string `json:"level"`
}

// HandleHealth GET /health
func (h *CoursesHandler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httputil.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	httputil.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// HandleCoursesCollection handles:
// - GET /api/v1/courses
// - POST /api/v1/courses
func (h *CoursesHandler) HandleCoursesCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleListCourses(w, r)
	case http.MethodPost:
		h.handleCreateCourse(w, r)
	default:
		httputil.Error(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

// HandleCourseItem handles:
// - GET /api/v1/courses/{id}
// - PATCH /api/v1/courses/{id}
// - DELETE /api/v1/courses/{id}
func (h *CoursesHandler) HandleCourseItem(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/courses/")
	if id == "" || strings.Contains(id, "/") {
		httputil.Error(w, http.StatusNotFound, "not found")
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.handleGetCourse(w, r, id)
	case http.MethodPatch:
		h.handleUpdateCourse(w, r, id)
	case http.MethodDelete:
		h.handleDeleteCourse(w, r, id)
	default:
		httputil.Error(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *CoursesHandler) handleListCourses(w http.ResponseWriter, r *http.Request) {
	if h.List == nil {
		httputil.Error(w, http.StatusInternalServerError, "list use case not configured")
		return
	}

	q := r.URL.Query()
	filter := usecase.ListCoursesFilter{
		Subject: q.Get("subject"),
		Level:   q.Get("level"),
		Limit:   20,
		Offset:  0,
	}
	if v := q.Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			filter.Limit = n
		}
	}
	if v := q.Get("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			filter.Offset = n
		}
	}

	courses, err := h.List.List(filter)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, "failed to list courses")
		return
	}

	resp := make([]CourseResponse, 0, len(courses))
	for _, c := range courses {
		resp = append(resp, courseToResponse(&c))
	}
	httputil.JSON(w, http.StatusOK, resp)
}

func (h *CoursesHandler) handleGetCourse(w http.ResponseWriter, r *http.Request, id string) {
	if h.Get == nil {
		httputil.Error(w, http.StatusInternalServerError, "get use case not configured")
		return
	}

	course, err := h.Get.Get(id)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, "failed to get course")
		return
	}
	if course == nil {
		httputil.Error(w, http.StatusNotFound, "course not found")
		return
	}
	httputil.JSON(w, http.StatusOK, courseToResponse(course))
}

func (h *CoursesHandler) handleCreateCourse(w http.ResponseWriter, r *http.Request) {
	if h.Create == nil {
		httputil.Error(w, http.StatusInternalServerError, "create use case not configured")
		return
	}

	userID := auth.UserIDFromContext(r.Context())
	if userID == "" {
		httputil.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req CreateCourseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid body")
		return
	}
	if req.Title == "" || req.Description == "" || req.Subject == "" || req.Level == "" {
		httputil.Error(w, http.StatusBadRequest, "title, description, subject, level required")
		return
	}

	in := usecase.CreateCourseInput{
		Title:       req.Title,
		Description: req.Description,
		Subject:     req.Subject,
		Level:       req.Level,
		OwnerUserID: userID,
	}
	course, err := h.Create.Create(in)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, "failed to create course")
		return
	}
	httputil.JSON(w, http.StatusCreated, courseToResponse(course))
}

func (h *CoursesHandler) handleUpdateCourse(w http.ResponseWriter, r *http.Request, id string) {
	if h.Update == nil {
		httputil.Error(w, http.StatusInternalServerError, "update use case not configured")
		return
	}

	userID := auth.UserIDFromContext(r.Context())
	if userID == "" {
		httputil.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req UpdateCourseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid body")
		return
	}

	in := usecase.UpdateCourseInput{
		ID:          id,
		Title:       req.Title,
		Description: req.Description,
		Subject:     req.Subject,
		Level:       req.Level,
	}

	course, err := h.Update.Update(in)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, "failed to update course")
		return
	}
	if course == nil {
		httputil.Error(w, http.StatusNotFound, "course not found")
		return
	}
	httputil.JSON(w, http.StatusOK, courseToResponse(course))
}

func (h *CoursesHandler) handleDeleteCourse(w http.ResponseWriter, r *http.Request, id string) {
	if h.Delete == nil {
		httputil.Error(w, http.StatusInternalServerError, "delete use case not configured")
		return
	}

	userID := auth.UserIDFromContext(r.Context())
	if userID == "" {
		httputil.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if err := h.Delete.Delete(id); err != nil {
		httputil.Error(w, http.StatusInternalServerError, "failed to delete course")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func courseToResponse(c *domain.Course) CourseResponse {
	return CourseResponse{
		ID:          c.ID,
		Title:       c.Title,
		Description: c.Description,
		Subject:     c.Subject,
		Level:       c.Level,
		OwnerUserID: c.OwnerUserID,
	}
}
