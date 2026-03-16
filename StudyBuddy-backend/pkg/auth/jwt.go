package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("invalid token")
)

// Claims holds JWT claims (access or refresh).
type Claims struct {
	jwt.RegisteredClaims
	UserID    string `json:"uid"`
	Email     string `json:"email"`
	IsRefresh bool   `json:"refresh,omitempty"`
}

// Config for JWT issuance and validation.
type Config struct {
	Secret     []byte
	Issuer     string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

// IssuePair returns access and refresh tokens for the user.
func IssuePair(cfg Config, userID, email string) (access, refresh string, expAt time.Time, err error) {
	now := time.Now()
	expAt = now.Add(cfg.AccessTTL)
	accessClaims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			Issuer:    cfg.Issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expAt),
			ID:        uuid.New().String(),
		},
		UserID:    userID,
		Email:     email,
		IsRefresh: false,
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	access, err = accessToken.SignedString(cfg.Secret)
	if err != nil {
		return "", "", time.Time{}, err
	}
	refreshExp := now.Add(cfg.RefreshTTL)
	refreshClaims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			Issuer:    cfg.Issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(refreshExp),
			ID:        uuid.New().String(),
		},
		UserID:    userID,
		Email:     email,
		IsRefresh: true,
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refresh, err = refreshToken.SignedString(cfg.Secret)
	if err != nil {
		return "", "", time.Time{}, err
	}
	return access, refresh, expAt, nil
}

// ValidateAccess parses the token and returns claims if it's a valid access token.
func ValidateAccess(secret []byte, tokenString string) (*Claims, error) {
	tok, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(*jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return nil, ErrInvalidToken
	}
	claims, ok := tok.Claims.(*Claims)
	if !ok || !tok.Valid || claims.IsRefresh {
		return nil, ErrInvalidToken
	}
	return claims, nil
}

// ValidateRefresh parses the token and returns claims if it's a valid refresh token.
func ValidateRefresh(secret []byte, tokenString string) (*Claims, error) {
	tok, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(*jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return nil, ErrInvalidToken
	}
	claims, ok := tok.Claims.(*Claims)
	if !ok || !tok.Valid || !claims.IsRefresh {
		return nil, ErrInvalidToken
	}
	return claims, nil
}
