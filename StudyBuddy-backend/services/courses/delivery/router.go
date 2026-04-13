package delivery

import (
	"net/http"
	"studybuddy/backend/pkg/auth"
)

// NewRouter returns the courses service HTTP router.
func NewRouter(h *CoursesHandler, jwtSecret []byte) http.Handler {
	protect := auth.Middleware(jwtSecret)
	mux := http.NewServeMux()

	// Health check.
	mux.HandleFunc("/health", h.HandleHealth)

	// Courses API.
	// Method is checked inside handlers to keep router simple (Go < 1.22 style).
	mux.Handle("/api/v1/courses", protect(http.HandlerFunc(h.HandleCoursesCollection)))
	mux.Handle("/api/v1/courses/", protect(http.HandlerFunc(h.HandleCourseItem)))

	return mux
}
