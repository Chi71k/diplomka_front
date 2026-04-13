package domain

import "time"

// Course represents a study course in the catalog.
type Course struct {
	ID          string
	Title       string
	Description string
	Subject     string
	Level       string
	OwnerUserID string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
