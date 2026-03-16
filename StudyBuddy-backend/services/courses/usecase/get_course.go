package usecase

import "studybuddy/backend/services/courses/domain"

// GetCourse defines the use case for retrieving a single course.
type GetCourse interface {
	Get(id string) (*domain.Course, error)
}
