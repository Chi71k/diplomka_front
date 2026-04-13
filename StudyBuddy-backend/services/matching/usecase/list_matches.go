package usecase

import (
	"fmt"
	"studybuddy/backend/services/matching/domain"
)

type ListMatchesInput struct {
	UserID string
	Status domain.MatchStatus // empty = all
	Limit  int
	Offset int
}

type ListMatches interface {
	List(in ListMatchesInput) ([]domain.Match, error)
}

type listMatches struct {
	repo MatchRepository
}

func NewListMatches(repo MatchRepository) ListMatches {
	return &listMatches{repo: repo}
}

func (uc *listMatches) List(in ListMatchesInput) ([]domain.Match, error) {
	limit := in.Limit
	if limit <= 0 {
		limit = 20
	}
	offset := in.Offset
	if offset < 0 {
		offset = 0
	}

	matches, err := uc.repo.ListForUser(in.UserID, ListMatchesFilter{
		Status: in.Status,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("list matches: %w", err)
	}
	return matches, nil
}
