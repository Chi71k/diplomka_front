package usecase

import "studybuddy/backend/services/users/domain"

// ProfileRepository is the port for profile persistence.
type ProfileRepository interface {
	GetByUserID(userID string) (*domain.Profile, error)
	Upsert(profile *domain.Profile) error
	// DeleteByUserID performs logical deletion for the user (e.g. deactivate).
	DeleteByUserID(userID string) error
}
