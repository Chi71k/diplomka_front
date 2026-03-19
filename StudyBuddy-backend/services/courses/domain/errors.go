package domain

import "errors"

var (
	ErrCourseNotFound = errors.New("course not found")
	ErrForbidden      = errors.New("forbidden: you do not own this course")
)
