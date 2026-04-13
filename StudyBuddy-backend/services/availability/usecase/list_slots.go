package usecase

import "studybuddy/backend/services/availability/domain"

type ListSlots interface {
	ListSlots(userID string) ([]domain.Slot, error)
}

type listSlots struct {
	repo SlotRepository
}

func NewListSlots(repo SlotRepository) ListSlots {
	return &listSlots{repo: repo}
}

func (l *listSlots) ListSlots(userID string) ([]domain.Slot, error) {
	return l.repo.ListForUser(userID)
}
