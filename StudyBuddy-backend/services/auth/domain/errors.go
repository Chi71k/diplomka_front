package domain

import "errors"

var (
	ErrEmailExists  = errors.New("email already registered")
	ErrInvalidCreds = errors.New("invalid credentials")
	ErrUserInactive = errors.New("user is inactive")
)
