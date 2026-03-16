package delivery

import "net/http"

// NewRouter returns the courses service HTTP router.
func NewRouter(h *CoursesHandler) http.Handler {
	mux := http.NewServeMux()

	// Health check.
	mux.HandleFunc("/health", h.HandleHealth)

	// Courses API.
	// Method is checked inside handlers to keep router simple (Go < 1.22 style).
	mux.HandleFunc("/api/v1/courses", h.HandleCoursesCollection)
	mux.HandleFunc("/api/v1/courses/", h.HandleCourseItem)

	return mux
}
