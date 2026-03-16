package usecase

import "studybuddy/backend/services/auth/domain"

// UserRepository is the port for user persistence (Auth service owns credentials).
type UserRepository interface {
	Create(user *domain.User) error
	GetByEmail(email string) (*domain.User, error)
	GetByID(id string) (*domain.User, error)
}

// PasswordHasher hashes and compares passwords.
type PasswordHasher interface {
	Hash(plain string) (string, error)
	Compare(hash, plain string) bool
}

// JWTIssuer issues token pairs (abstracts pkg/auth for testability).
type JWTIssuer interface {
	IssuePair(userID, email string) (access, refresh string, expAtUnix int64, err error)
}
