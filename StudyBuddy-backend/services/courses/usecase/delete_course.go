package usecase

// DeleteCourse defines the use case for deleting a course.
type DeleteCourse interface {
	Delete(id string) error
}
