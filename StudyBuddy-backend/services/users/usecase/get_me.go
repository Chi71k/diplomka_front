package usecase

import (
	"studybuddy/backend/services/users/domain"
)

// GetMe returns the profile for the given user ID (from JWT).
type GetMe interface {
	GetMe(userID string) (*domain.Profile, error)
}

type getMe struct {
	repo ProfileRepository
}

// NewGetMe creates the GetMe use case.
func NewGetMe(repo ProfileRepository) GetMe {
	return &getMe{repo: repo}
}

func (u *getMe) GetMe(userID string) (*domain.Profile, error) {
	return u.repo.GetByUserID(userID)
}
