package usecase

import "studybuddy/backend/services/courses/domain"

// UpdateCourseInput is the input for updating a course.
type UpdateCourseInput struct {
	ID          string
	Title       *string
	Description *string
	Subject     *string
	Level       *string
}

// UpdateCourse defines the use case for updating a course.
type UpdateCourse interface {
	Update(input UpdateCourseInput) (*domain.Course, error)
}
