package usecase

import "studybuddy/backend/services/courses/domain"

// CourseRepository is the port for course persistence.
type CourseRepository interface {
	Create(course *domain.Course) error
	Update(course *domain.Course) error
	Delete(id string) error
	GetByID(id string) (*domain.Course, error)
	List(filter ListCoursesFilter) ([]domain.Course, error)
}

// ListCoursesFilter defines basic filters for listing courses.
type ListCoursesFilter struct {
	Subject string
	Level   string
	Limit   int
	Offset  int
}
