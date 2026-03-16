package usecase

import "studybuddy/backend/services/courses/domain"

// Service aggregates all course use cases.
type Service struct {
	repo CourseRepository
}

// NewService creates a new Service.
func NewService(repo CourseRepository) *Service {
	return &Service{repo: repo}
}

// Ensure Service implements interfaces.
var (
	_ CreateCourse = (*Service)(nil)
	_ GetCourse    = (*Service)(nil)
	_ ListCourses  = (*Service)(nil)
	_ UpdateCourse = (*Service)(nil)
	_ DeleteCourse = (*Service)(nil)
)

func (s *Service) Create(input CreateCourseInput) (*domain.Course, error) {
	c := &domain.Course{
		Title:       input.Title,
		Description: input.Description,
		Subject:     input.Subject,
		Level:       input.Level,
		OwnerUserID: input.OwnerUserID,
	}
	if err := s.repo.Create(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) Get(id string) (*domain.Course, error) {
	return s.repo.GetByID(id)
}

func (s *Service) List(filter ListCoursesFilter) ([]domain.Course, error) {
	return s.repo.List(filter)
}

func (s *Service) Update(input UpdateCourseInput) (*domain.Course, error) {
	existing, err := s.repo.GetByID(input.ID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, nil
	}
	if input.Title != nil {
		existing.Title = *input.Title
	}
	if input.Description != nil {
		existing.Description = *input.Description
	}
	if input.Subject != nil {
		existing.Subject = *input.Subject
	}
	if input.Level != nil {
		existing.Level = *input.Level
	}
	if err := s.repo.Update(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *Service) Delete(id string) error {
	return s.repo.Delete(id)
}
