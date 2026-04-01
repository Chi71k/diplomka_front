package usecase

import (
	"fmt"
	"studybuddy/backend/services/matching/domain"
)

type SendMatchRequestInput struct {
	RequesterID string
	ReceiverID  string
	Message     string // optional
}

type SendMatchRequest interface {
	Send(in SendMatchRequestInput) (*domain.Match, error)
}

type sendMatchRequest struct {
	repo MatchRepository
}

func NewSendMatchRequest(repo MatchRepository) SendMatchRequest {
	return &sendMatchRequest{repo: repo}
}

func (uc *sendMatchRequest) Send(in SendMatchRequestInput) (*domain.Match, error) {
	if in.RequesterID == in.ReceiverID {
		return nil, domain.ErrCannotMatchSelf
	}

	existing, err := uc.repo.GetBetween(in.RequesterID, in.ReceiverID)
	if err != nil {
		return nil, fmt.Errorf("check existing match: %w", err)
	}
	if existing != nil &&
		(existing.Status == domain.MatchStatusPending || existing.Status == domain.MatchStatusAccepted) {
		return nil, domain.ErrMatchAlreadyExists
	}

	m := &domain.Match{
		RequesterID: in.RequesterID,
		ReceiverID:  in.ReceiverID,
		Status:      domain.MatchStatusPending,
		Message:     in.Message,
	}
	if err := uc.repo.Create(m); err != nil {
		return nil, fmt.Errorf("create match: %w", err)
	}
	return m, nil
}
