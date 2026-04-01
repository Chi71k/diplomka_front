package usecase

import (
	"context"
	"studybuddy/backend/services/availability/domain"
)

type SlotRepository interface {
	Create(slot *domain.Slot) error
	ListForUser(userID string) ([]domain.Slot, error)
	GetByID(id string) (*domain.Slot, error)
	Delete(id string) error
	// DeleteAllForUser removes every slot belonging to a user
	// Example: user disconnects GCal and opts out to clear imported slots,
	// account deletion
	DeleteAllForUser(userID string) error

	// ReplaceForUser atomically replaces all slots for a user with the provided slice
	// It is used by the GCal import flow to keep local slots in sync
	ReplaceForUser(userID string, slots []domain.Slot) error

	// TODO: implement this method for Matching service later
	ListForUsers(userIDs []string) ([]domain.Slot, error)
}

type GCalRepository interface {
	GetConnection(userID string) (*domain.GCalConnection, error)
	UpsertConnection(connection *domain.GCalConnection) error
	DeleteConnection(userID string) error
}

// GCalProvider is the port for Google Calendar API
type GCalProvider interface {
	// ExchangeCode exchanges OAuth code for tokens
	ExchangeCode(ctx context.Context, code string) (*domain.GCalConnection, error)
	// RefreshToken refreshes an expired access token
	RefreshToken(ctx context.Context, conn *domain.GCalConnection) (*domain.GCalConnection, error)
	// ImportEvents fetches busy slots from GCal and converts to []domain.Slot
	ImportEvents(ctx context.Context, conn *domain.GCalConnection, userID string) ([]domain.Slot, error)
	// GetAuthURL returns the OAuth redirect the URL for the frontend
	GetAuthURL(state string) string
}
