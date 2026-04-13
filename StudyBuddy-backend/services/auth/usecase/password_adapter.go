package usecase

import "studybuddy/backend/pkg/password"

// PasswordAdapter adapts pkg/password to usecase.PasswordHasher.
type PasswordAdapter struct{}

// Hash implements PasswordHasher.
func (PasswordAdapter) Hash(plain string) (string, error) {
	return password.Hash(plain)
}

// Compare implements PasswordHasher.
func (PasswordAdapter) Compare(hash, plain string) bool {
	return password.Compare(hash, plain)
}
