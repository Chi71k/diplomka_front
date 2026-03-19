package usecase

type DeleteCourseInput struct {
	ID             string
	RequestingUser string // for ownership check (JWT)
}

// DeleteCourse defines the use case for deleting a course.
type DeleteCourse interface {
	Delete(input DeleteCourseInput) error
}
