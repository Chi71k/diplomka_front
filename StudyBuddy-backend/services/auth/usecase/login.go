package usecase

import (
	"studybuddy/backend/services/auth/domain"
)

// LoginInput for login.
type LoginInput struct {
	Email    string
	Password string
}

// LoginOutput returned after successful login.
type LoginOutput struct {
	UserID       string
	Email        string
	AccessToken  string
	RefreshToken string
	ExpiresAt    int64
}

// Login authenticates and returns tokens.
type Login interface {
	Login(in LoginInput) (LoginOutput, error)
}

type login struct {
	repo   UserRepository
	hasher PasswordHasher
	jwt    JWTIssuer
}

// NewLogin creates the Login use case.
func NewLogin(repo UserRepository, hasher PasswordHasher, jwt JWTIssuer) Login {
	return &login{repo: repo, hasher: hasher, jwt: jwt}
}

func (u *login) Login(in LoginInput) (LoginOutput, error) {
	user, err := u.repo.GetByEmail(in.Email)
	if err != nil || user == nil {
		return LoginOutput{}, domain.ErrInvalidCreds
	}
	if !user.IsActive {
		return LoginOutput{}, domain.ErrUserInactive
	}
	if !u.hasher.Compare(user.PasswordHash, in.Password) {
		return LoginOutput{}, domain.ErrInvalidCreds
	}
	access, refresh, expAt, err := u.jwt.IssuePair(user.ID, user.Email)
	if err != nil {
		return LoginOutput{}, err
	}
	return LoginOutput{
		UserID:       user.ID,
		Email:        user.Email,
		AccessToken:  access,
		RefreshToken: refresh,
		ExpiresAt:    expAt,
	}, nil
}
