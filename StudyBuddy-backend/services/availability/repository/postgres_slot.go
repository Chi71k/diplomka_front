package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"studybuddy/backend/services/availability/domain"
	"studybuddy/backend/services/availability/usecase"
)

type PgSlotRepository struct {
	pool *pgxpool.Pool
}

func NewPgSlotRepository(pool *pgxpool.Pool) usecase.SlotRepository {
	return &PgSlotRepository{pool: pool}
}

var _ usecase.SlotRepository = (*PgSlotRepository)(nil)

func (r *PgSlotRepository) Create(slot *domain.Slot) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	const q = `
INSERT INTO availability_slots (user_id, day_of_week, start_time, end_time, timezone)
VALUES ($1, $2, $3::time, $4::time, $5)
RETURNING id, created_at, updated_at;
`
	err := r.pool.QueryRow(ctx, q,
		slot.UserID,
		slot.DayOfWeek,
		slot.StartTime.Format("15:04:05"),
		slot.EndTime.Format("15:04:05"),
		slot.Timezone,
	).Scan(&slot.ID, &slot.CreatedAt, &slot.UpdatedAt)
	if err != nil {
		if isUniqueViolation(err) {
			return domain.ErrSlotConflict
		}
		return fmt.Errorf("insert slot: %w", err)
	}
	return nil
}

func (r *PgSlotRepository) ListForUser(userID string) ([]domain.Slot, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	const q = `
SELECT id, user_id, day_of_week, start_time, end_time, timezone, created_at, updated_at
FROM availability_slots
WHERE user_id = $1
ORDER BY day_of_week, start_time;
`
	rows, err := r.pool.Query(ctx, q, userID)
	if err != nil {
		return nil, fmt.Errorf("list slots for user: %w", err)
	}
	defer rows.Close()

	return scanSlots(rows)
}

func (r *PgSlotRepository) GetByID(id string) (*domain.Slot, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	const q = `
SELECT id, user_id, day_of_week, start_time, end_time, timezone, created_at, updated_at
FROM availability_slots
WHERE id = $1;
`
	var slot domain.Slot

	err := r.pool.QueryRow(ctx, q, id).Scan(
		&slot.ID,
		&slot.UserID,
		&slot.DayOfWeek,
		&slot.StartTime,
		&slot.EndTime,
		&slot.Timezone,
		&slot.CreatedAt,
		&slot.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get slot by id: %w", err)
	}

	return &slot, nil
}

func (r *PgSlotRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	const q = `DELETE FROM availability_slots WHERE id = $1;`

	if _, err := r.pool.Exec(ctx, q, id); err != nil {
		return fmt.Errorf("delete slot: %w", err)
	}
	return nil
}

func (r *PgSlotRepository) DeleteAllForUser(userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	const q = `DELETE FROM availability_slots WHERE user_id = $1;`

	if _, err := r.pool.Exec(ctx, q, userID); err != nil {
		return fmt.Errorf("delete all slots for user: %w", err)
	}
	return nil
}

func (r *PgSlotRepository) ReplaceForUser(userID string, slots []domain.Slot) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if _, err := tx.Exec(ctx, `DELETE FROM availability_slots WHERE user_id = $1`, userID); err != nil {
		return fmt.Errorf("delete existing slots in replace: %w", err)
	}

	const ins = `
INSERT INTO availability_slots (user_id, day_of_week, start_time, end_time, timezone)
VALUES ($1, $2, $3::time, $4::time, $5)
ON CONFLICT (user_id, day_of_week, start_time) DO NOTHING;
`
	for i := range slots {
		s := &slots[i]
		if err := s.Validate(); err != nil {
			continue // skipping invalid slots
		}
		if _, err := tx.Exec(ctx, ins,
			userID,
			s.DayOfWeek,
			s.StartTime.Format("15:04:05"),
			s.EndTime.Format("15:04:05"),
			s.Timezone,
		); err != nil {
			return fmt.Errorf("insert slot during replace (day=%d start=%s): %w",
				s.DayOfWeek, s.StartTime.Format("15:04"), err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit replace transaction: %w", err)
	}
	return nil
}

func (r *PgSlotRepository) ListForUsers(userIDs []string) ([]domain.Slot, error) {
	if len(userIDs) == 0 {
		return []domain.Slot{}, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	const q = `
SELECT id, user_id, day_of_week, start_time, end_time, timezone, created_at, updated_at
FROM availability_slots
WHERE user_id = ANY($1::uuid[])
ORDER BY user_id, day_of_week, start_time;
`
	rows, err := r.pool.Query(ctx, q, userIDs)
	if err != nil {
		return nil, fmt.Errorf("list slots for users: %w", err)
	}
	defer rows.Close()

	return scanSlots(rows)
}

// helpers
type pgRows interface {
	Next() bool
	Scan(dest ...any) error
	Err() error
	Close()
}

func scanSlots(rows pgRows) ([]domain.Slot, error) {
	var out []domain.Slot
	for rows.Next() {
		var slot domain.Slot
		if err := rows.Scan(
			&slot.ID,
			&slot.UserID,
			&slot.DayOfWeek,
			&slot.StartTime,
			&slot.EndTime,
			&slot.Timezone,
			&slot.CreatedAt,
			&slot.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan slot row: %w", err)
		}
		out = append(out, slot)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return out, nil
}

func isUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	type pgErr interface {
		SQLState() string
	}
	var pg pgErr
	if errors.As(err, &pg) {
		return pg.SQLState() == "23505"
	}
	return false
}
