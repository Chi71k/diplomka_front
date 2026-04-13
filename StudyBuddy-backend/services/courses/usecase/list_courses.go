package usecase

import "studybuddy/backend/services/courses/domain"

// ListCourses defines the use case for listing courses.
type ListCourses interface {
	List(filter ListCoursesFilter) ([]domain.Course, error)
}
