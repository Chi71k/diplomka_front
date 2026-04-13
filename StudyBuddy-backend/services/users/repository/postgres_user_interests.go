package repository

import (
	"context"
	"studybuddy/backend/services/users/domain"
	"studybuddy/backend/services/users/usecase"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PgUserInterestRepository struct {
	pool *pgxpool.Pool
}

func NewPgUserInterestRepository(pool *pgxpool.Pool) usecase.UserInterestRepository {
	return &PgUserInterestRepository{pool: pool}
}

func (r *PgUserInterestRepository) ListForUser(userID string) ([]domain.Interest, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	const q = `
SELECT i.id, i.name
FROM user_interests ui
JOIN interests i ON i.id = ui.interest_id
WHERE ui.user_id = $1
ORDER BY i.name;
`
	rows, err := r.pool.Query(ctx, q, userID)
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

func (r *PgUserInterestRepository) ReplaceForUser(userID string, interestIDs []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	if _, err := tx.Exec(ctx, "DELETE FROM user_interests WHERE user_id = $1", userID); err != nil {
		return err
	}

	if len(interestIDs) > 0 {
		const ins = `
INSERT INTO user_interests (user_id, interest_id)
SELECT $1, unnest($2::text[])::uuid
ON CONFLICT DO NOTHING;
`
		if _, err := tx.Exec(ctx, ins, userID, interestIDs); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
