package usecase

import (
	"fmt"
	"studybuddy/backend/services/matching/domain"
)

type CancelMatchInput struct {
	MatchID     string
	RequesterID string
}

type CancelMatch interface {
	Cancel(in CancelMatchInput) error
}

type cancelMatch struct {
	repo MatchRepository
}

func NewCancelMatch(repo MatchRepository) CancelMatch {
	return &cancelMatch{repo: repo}
}

func (uc *cancelMatch) Cancel(in CancelMatchInput) error {
	m, err := uc.repo.GetByID(in.MatchID)
	if err != nil {
		return fmt.Errorf("get match: %w", err)
	}
	if m == nil {
		return domain.ErrMatchNotFound
	}
	// Only the requester may cancel.
	if m.RequesterID != in.RequesterID {
		return domain.ErrForbidden
	}
	if m.Status != domain.MatchStatusPending {
		return domain.ErrInvalidStatusChange
	}

	if err := uc.repo.UpdateStatus(m.ID, domain.MatchStatusCanceled); err != nil {
		return fmt.Errorf("cancel match: %w", err)
	}
	return nil
}
