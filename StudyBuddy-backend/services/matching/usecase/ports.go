package usecase

import (
	"studybuddy/backend/services/matching/domain"
)

type MatchRepository interface {
	Create(m *domain.Match) error
	GetByID(id string) (*domain.Match, error)
	// GetBetween is used for pending/accepted match between the two users
	GetBetween(userA, userB string) (*domain.Match, error)
	UpdateStatus(id string, status domain.MatchStatus) error
	ListForUser(userID string, filter ListMatchesFilter) ([]domain.Match, error)
}

// ListMatchesFilter narrows ListForUser results
type ListMatchesFilter struct {
	// Status filters by a specific status. Empty string means "all"
	Status domain.MatchStatus
	Limit  int
	Offset int
}

// ProfileClient fetches the lightweight profile data from the Users service
type ProfileClient interface {
	GetProfile(userID string) (*ProfileData, error)
	GetInterestIDs(userID string) ([]string, error)
}

// ProfileData is the minimal profile needed for candidate enrichment
type ProfileData struct {
	UserID    string
	FirstName string
	LastName  string
	Bio       string
	AvatarURL string
}

// SlotClient fetches the availability slots from the Availability service
type SlotClient interface {
	ListForUser(userID string) ([]SlotData, error)
	ListForUsers(userIDs []string) ([]SlotData, error)
}

// SlotData is the minimal slot needed for overlap scoring
type SlotData struct {
	ID        string
	UserID    string
	DayOfWeek int
	StartTime string
	EndTime   string
	Timezone  string
}

// CourseClient fetches courses the user owns/enrolled in
// For MVP, a user "has" courses they created
type CourseClient interface {
	ListCourseIDsForUser(userID string) ([]string, error)
}

// CandidateStore provides a pre-filtered pool of candidate user IDs
// Lists all active users except the requester and those already matched
type CandidateStore interface {
	ListCandidateIDs(requesterID string, excludeIDs []string) ([]string, error)
}
