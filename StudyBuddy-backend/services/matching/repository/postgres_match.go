package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"studybuddy/backend/services/matching/domain"
	"studybuddy/backend/services/matching/usecase"
)

type PgMatchRepository struct {
	pool *pgxpool.Pool
}

func NewPgMatchRepository(pool *pgxpool.Pool) usecase.MatchRepository {
	return &PgMatchRepository{pool: pool}
}

func (r *PgMatchRepository) Create(m *domain.Match) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	const q = `
INSERT INTO matches (requester_id, receiver_id, status, message)
VALUES ($1, $2, $3, $4)
RETURNING id, created_at, updated_at;
`
	return r.pool.QueryRow(ctx, q,
		m.RequesterID,
		m.ReceiverID,
		string(m.Status),
		m.Message,
	).Scan(&m.ID, &m.CreatedAt, &m.UpdatedAt)
}

func (r *PgMatchRepository) GetByID(id string) (*domain.Match, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	const q = `
SELECT id, requester_id, receiver_id, status, message, created_at, updated_at
FROM matches
WHERE id = $1;
`
	m, err := scanMatch(r.pool.QueryRow(ctx, q, id))
	if errors.Is(err, pgx.ErrNoRows) || isInvalidUUID(err) {
		return nil, nil
	}
	return m, err
}

func (r *PgMatchRepository) GetBetween(userA, userB string) (*domain.Match, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Returns any match (regardless of direction) that is pending or accepted.
	const q = `
SELECT id, requester_id, receiver_id, status, message, created_at, updated_at
FROM matches
WHERE (
    (requester_id = $1 AND receiver_id = $2)
 OR (requester_id = $2 AND receiver_id = $1)
)
  AND status IN ('pending', 'accepted')
LIMIT 1;
`
	m, err := scanMatch(r.pool.QueryRow(ctx, q, userA, userB))
	if errors.Is(err, pgx.ErrNoRows) || isInvalidUUID(err) {
		return nil, nil
	}
	return m, err
}

func (r *PgMatchRepository) UpdateStatus(id string, status domain.MatchStatus) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	const q = `
UPDATE matches
SET status = $2, updated_at = now()
WHERE id = $1;
`
	ct, err := r.pool.Exec(ctx, q, id, status)
	if err != nil {
		if isInvalidUUID(err) {
			return domain.ErrMatchNotFound
		}
	}
	if ct.RowsAffected() == 0 {
		return domain.ErrMatchNotFound
	}
	return nil
}

func (r *PgMatchRepository) ListForUser(userID string, f usecase.ListMatchesFilter) ([]domain.Match, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	args := []any{userID}
	where := ""
	if f.Status != "" {
		args = append(args, string(f.Status))
		where = fmt.Sprintf(" AND status = $%d", len(args))
	}
	limit := f.Limit
	if limit <= 0 {
		limit = 20
	}
	offset := f.Offset
	args = append(args, limit)
	args = append(args, offset)

	q := fmt.Sprintf(`
SELECT id, requester_id, receiver_id, status, message, created_at, updated_at
FROM matches
WHERE (requester_id = $1 OR receiver_id = $1)
%s
ORDER BY created_at DESC
LIMIT $%d OFFSET $%d;
`, where, len(args)-1, len(args))

	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("list matches: %w", err)
	}
	defer rows.Close()

	var out []domain.Match
	for rows.Next() {
		m, err := scanMatch(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *m)
	}
	return out, rows.Err()
}

// shared scanner
type rowScanner interface {
	Scan(dest ...any) error
}

func scanMatch(row rowScanner) (*domain.Match, error) {
	var m domain.Match
	var status string
	err := row.Scan(
		&m.ID,
		&m.RequesterID,
		&m.ReceiverID,
		&status,
		&m.Message,
		&m.CreatedAt,
		&m.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	m.Status = domain.MatchStatus(status)
	return &m, nil
}

// helpers
func isInvalidUUID(err error) bool {
	return strings.Contains(err.Error(), "22PO2") ||
		strings.Contains(err.Error(), "invalid input syntax for type uuid")
}
