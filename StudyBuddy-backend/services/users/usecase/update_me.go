package usecase

import (
	"studybuddy/backend/services/users/domain"
)

// UpdateMeInput for partial profile update.
type UpdateMeInput struct {
	UserID    string
	FirstName *string
	LastName  *string
	Bio       *string
	AvatarURL *string
}

// UpdateMe updates the profile for the given user ID.
type UpdateMe interface {
	UpdateMe(in UpdateMeInput) (*domain.Profile, error)
}

type updateMe struct {
	repo ProfileRepository
}

// NewUpdateMe creates the UpdateMe use case.
func NewUpdateMe(repo ProfileRepository) UpdateMe {
	return &updateMe{repo: repo}
}

func (u *updateMe) UpdateMe(in UpdateMeInput) (*domain.Profile, error) {
	existing, err := u.repo.GetByUserID(in.UserID)
	if err != nil {
		return nil, err
	}
	profile := &domain.Profile{UserID: in.UserID}
	if existing != nil {
		profile.FirstName = existing.FirstName
		profile.LastName = existing.LastName
		profile.Bio = existing.Bio
		profile.AvatarURL = existing.AvatarURL
		profile.Email = existing.Email
		profile.CreatedAt = existing.CreatedAt
	}
	if in.FirstName != nil {
		profile.FirstName = *in.FirstName
	}
	if in.LastName != nil {
		profile.LastName = *in.LastName
	}
	if in.Bio != nil {
		profile.Bio = *in.Bio
	}
	if in.AvatarURL != nil {
		profile.AvatarURL = *in.AvatarURL
	}
	if err := u.repo.Upsert(profile); err != nil {
		return nil, err
	}
	return u.repo.GetByUserID(in.UserID)
}
