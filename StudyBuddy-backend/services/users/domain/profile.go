package domain

import "time"

// Profile is the user profile (owned by Users service).
type Profile struct {
	UserID    string
	Email     string // optional, for display; may be synced from Auth
	FirstName string
	LastName  string
	Bio       string
	AvatarURL string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Interest for interests selection.
type Interest struct {
	ID   string
	Name string
}
