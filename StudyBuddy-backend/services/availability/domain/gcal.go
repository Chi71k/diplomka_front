package domain

import "time"

type GCalConnection struct {
	UserID       string
	AccessToken  string // should be encrypted
	RefreshToken string
	TokenExpiry  time.Time
	CalendarID   string
	SyncEnabled  bool
	LastSyncedAt *time.Time
}
