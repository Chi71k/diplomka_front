package repository

import (
	"context"
	"errors"
	"time"

	"studybuddy/backend/services/users/domain"
	"studybuddy/backend/services/users/usecase"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PgProfileRepository implements ProfileRepository using PostgreSQL users table.
type PgProfileRepository struct {
	pool *pgxpool.Pool
}

// NewPgProfileRepository creates a new PgProfileRepository.
func NewPgProfileRepository(pool *pgxpool.Pool) usecase.ProfileRepository {
	return &PgProfileRepository{pool: pool}
}

func (r *PgProfileRepository) GetByUserID(userID string) (*domain.Profile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	const q = `
SELECT id, email, first_name, last_name, bio, avatar_url, created_at, updated_at
FROM users
WHERE id = $1;
`
	var p domain.Profile
	err := r.pool.QueryRow(ctx, q, userID).Scan(
		&p.UserID,
		&p.Email,
		&p.FirstName,
		&p.LastName,
		&p.Bio,
		&p.AvatarURL,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

func (r *PgProfileRepository) Upsert(profile *domain.Profile) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Users service must not create new users; it only updates existing profiles.
	// We only update profile fields and leave credentials (password_hash) untouched.
	const q = `
UPDATE users
SET first_name = $2,
    last_name  = $3,
    bio        = $4,
    avatar_url = $5,
    updated_at = now()
WHERE id = $1;
`
	_, err := r.pool.Exec(ctx, q,
		profile.UserID,
		profile.FirstName,
		profile.LastName,
		profile.Bio,
		profile.AvatarURL,
	)
	return err
}

// DeleteByUserID performs logical deletion: mark user as inactive.
// Auth service will then reject logins for this user.
func (r *PgProfileRepository) DeleteByUserID(userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	const q = `
UPDATE users
SET is_active = false,
    updated_at = now()
WHERE id = $1;
`
	_, err := r.pool.Exec(ctx, q, userID)
	return err
}
