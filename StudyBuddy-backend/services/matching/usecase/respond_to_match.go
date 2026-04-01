package usecase

import (
	"fmt"
	"studybuddy/backend/services/matching/domain"
)

type RespondToMatchInput struct {
	MatchID     string
	ResponderID string
	Accept      bool // true = accept, false = decline
}

type RespondToMatch interface {
	Respond(in RespondToMatchInput) (*domain.Match, error)
}

type respondToMatch struct {
	repo MatchRepository
}

func NewRespondToMatch(repo MatchRepository) RespondToMatch {
	return &respondToMatch{repo: repo}
}

func (uc *respondToMatch) Respond(in RespondToMatchInput) (*domain.Match, error) {
	m, err := uc.repo.GetByID(in.MatchID)
	if err != nil {
		return nil, fmt.Errorf("get match: %w", err)
	}
	if m == nil {
		return nil, domain.ErrMatchNotFound
	}
	// Only the receiver may respond.
	if m.ReceiverID != in.ResponderID {
		return nil, domain.ErrForbidden
	}
	if m.Status != domain.MatchStatusPending {
		return nil, domain.ErrInvalidStatusChange
	}

	newStatus := domain.MatchStatusDeclined
	if in.Accept {
		newStatus = domain.MatchStatusAccepted
	}

	if err := uc.repo.UpdateStatus(m.ID, newStatus); err != nil {
		return nil, fmt.Errorf("update match status: %w", err)
	}
	m.Status = newStatus
	return m, nil
}
