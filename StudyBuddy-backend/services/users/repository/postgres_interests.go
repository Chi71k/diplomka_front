package repository

import (
	"context"
	"time"

	"studybuddy/backend/services/users/domain"
	"studybuddy/backend/services/users/usecase"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PgInterestRepository implements InterestRepository using the interests table.
type PgInterestRepository struct {
	pool *pgxpool.Pool
}

// NewPgInterestRepository creates a new PgInterestRepository.
func NewPgInterestRepository(pool *pgxpool.Pool) usecase.InterestRepository {
	return &PgInterestRepository{pool: pool}
}

func (r *PgInterestRepository) ListAll() ([]domain.Interest, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	const q = `
SELECT id, name
FROM interests
ORDER BY name;
`
	rows, err := r.pool.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.Interest
	for rows.Next() {
		var i domain.Interest
		if err := rows.Scan(&i.ID, &i.Name); err != nil {
			return nil, err
		}
		out = append(out, i)
	}
	return out, rows.Err()
}

func (r *PgInterestRepository) GetByIDs(ids []string) ([]domain.Interest, error) {
	if len(ids) == 0 {
		return []domain.Interest{}, nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	const q = `
SELECT id, name
FROM interests
WHERE id = ANY($1::uuid[]);
`
	rows, err := r.pool.Query(ctx, q, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.Interest
	for rows.Next() {
		var i domain.Interest
		if err := rows.Scan(&i.ID, &i.Name); err != nil {
			return nil, err
		}
		out = append(out, i)
	}
	return out, rows.Err()
}
