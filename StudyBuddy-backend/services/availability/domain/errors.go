package domain

import "errors"

var (
	// OAuth state errors
	ErrInvalidOAuthState = errors.New("invalid or tampered oauth state")
	ErrOAuthStateExpired = errors.New("OAuth state has expired, please try again")

	// GCal connection errors
	ErrGCalNotConnected = errors.New("google calendar is not connected for this user")
	ErrGCalSyncDisabled = errors.New("google calendar sync is disabled for this user")
)
