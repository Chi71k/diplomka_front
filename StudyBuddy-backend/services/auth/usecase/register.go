package usecase

import (
	"studybuddy/backend/services/auth/domain"
)

// RegisterInput for registration.
type RegisterInput struct {
	Email     string
	Password  string
	FirstName string
	LastName  string
}

// RegisterOutput returned after successful registration.
type RegisterOutput struct {
	UserID       string
	Email        string
	AccessToken  string
	RefreshToken string
	ExpiresAt    int64 // unix seconds
}

// Register registers a new user and returns tokens.
type Register interface {
	Register(in RegisterInput) (RegisterOutput, error)
}

type register struct {
	repo   UserRepository
	hasher PasswordHasher
	jwt    JWTIssuer
}

// NewRegister creates the Register use case.
func NewRegister(repo UserRepository, hasher PasswordHasher, jwt JWTIssuer) Register {
	return &register{repo: repo, hasher: hasher, jwt: jwt}
}

func (u *register) Register(in RegisterInput) (RegisterOutput, error) {
	existing, _ := u.repo.GetByEmail(in.Email)
	if existing != nil {
		return RegisterOutput{}, domain.ErrEmailExists
	}
	hash, err := u.hasher.Hash(in.Password)
	if err != nil {
		return RegisterOutput{}, err
	}
	user := &domain.User{
		Email:        in.Email,
		PasswordHash: hash,
		FirstName:    in.FirstName,
		LastName:     in.LastName,
		IsActive:     true,
	}
	if err := u.repo.Create(user); err != nil {
		return RegisterOutput{}, err
	}
	access, refresh, expAt, err := u.jwt.IssuePair(user.ID, user.Email)
	if err != nil {
		return RegisterOutput{}, err
	}
	return RegisterOutput{
		UserID:       user.ID,
		Email:        user.Email,
		AccessToken:  access,
		RefreshToken: refresh,
		ExpiresAt:    expAt,
	}, nil
}
