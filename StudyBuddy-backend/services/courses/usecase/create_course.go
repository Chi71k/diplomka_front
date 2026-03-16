package usecase

import "studybuddy/backend/services/courses/domain"

// CreateCourseInput is the input for creating a course.
type CreateCourseInput struct {
	Title       string
	Description string
	Subject     string
	Level       string
	OwnerUserID string
}

// CreateCourse defines the use case for creating a course.
type CreateCourse interface {
	Create(input CreateCourseInput) (*domain.Course, error)
}
