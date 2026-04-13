package usecase

import "studybuddy/backend/services/users/domain"

// InterestRepository for listing and resolving interests.
type InterestRepository interface {
	ListAll() ([]domain.Interest, error)
	GetByIDs(ids []string) ([]domain.Interest, error)
}

type UserInterestRepository interface {
	ListForUser(userID string) ([]domain.Interest, error)
	ReplaceForUser(userID string, interestIDs []string) error
}
