package usecase

import (
	"studybuddy/backend/pkg/auth"
)

// JWTAdapter adapts pkg/auth to usecase.JWTIssuer.
type JWTAdapter struct {
	Config auth.Config
}

// IssuePair implements JWTIssuer.
func (a JWTAdapter) IssuePair(userID, email string) (access, refresh string, expAtUnix int64, err error) {
	access, refresh, expAt, err := auth.IssuePair(a.Config, userID, email)
	if err != nil {
		return "", "", 0, err
	}
	return access, refresh, expAt.Unix(), nil
}
