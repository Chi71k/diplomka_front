package domain

import "time"

// User represents an authenticated user (credentials + minimal profile for register).
type User struct {
	ID           string
	Email        string
	PasswordHash string
	FirstName    string
	LastName     string
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
