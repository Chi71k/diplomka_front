package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"studybuddy/backend/services/matching/usecase"
)

type PgCandidateStore struct {
	pool *pgxpool.Pool
}

func NewPgCandidateStore(pool *pgxpool.Pool) usecase.CandidateStore {
	return &PgCandidateStore{pool: pool}
}

func (r *PgCandidateStore) ListCandidateIDs(requesterID string, excludeIDs []string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Exclude the requester themselves and all already-matched users.
	excluded := append([]string{requesterID}, excludeIDs...)

	const q = `
SELECT id::text
FROM users
WHERE is_active = true
  AND id != ALL($1::uuid[])
ORDER BY created_at DESC
LIMIT 200;
`
	rows, err := r.pool.Query(ctx, q, excluded)
	if err != nil {
		return nil, fmt.Errorf("list candidate ids: %w", err)
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}
