package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"studybuddy/backend/services/matching/usecase"
)

// Profile client
type PgProfileClient struct {
	pool *pgxpool.Pool
}

func NewPgProfileClient(pool *pgxpool.Pool) usecase.ProfileClient {
	return &PgProfileClient{pool: pool}
}

func (c *PgProfileClient) GetProfile(userID string) (*usecase.ProfileData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	const q = `
SELECT id::text, first_name, last_name,
       COALESCE(bio, ''), COALESCE(avatar_url, '')
FROM users
WHERE id = $1 AND is_active = true;
`
	var p usecase.ProfileData
	err := c.pool.QueryRow(ctx, q, userID).Scan(
		&p.UserID, &p.FirstName, &p.LastName, &p.Bio, &p.AvatarURL,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get profile: %w", err)
	}
	return &p, nil
}

func (c *PgProfileClient) GetInterestIDs(userID string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	const q = `
SELECT interest_id::text
FROM user_interests
WHERE user_id = $1;
`
	rows, err := c.pool.Query(ctx, q, userID)
	if err != nil {
		return nil, fmt.Errorf("get interest ids: %w", err)
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

// Slot client
type PgSlotClient struct {
	pool *pgxpool.Pool
}

func NewPgSlotClient(pool *pgxpool.Pool) usecase.SlotClient {
	return &PgSlotClient{pool: pool}
}

func (c *PgSlotClient) ListForUser(userID string) ([]usecase.SlotData, error) {
	return c.querySlots(
		`SELECT id::text, user_id::text, day_of_week, start_time::text, end_time::text, timezone
		 FROM availability_slots WHERE user_id = $1`,
		userID,
	)
}

func (c *PgSlotClient) ListForUsers(userIDs []string) ([]usecase.SlotData, error) {
	if len(userIDs) == 0 {
		return nil, nil
	}
	return c.querySlots(
		`SELECT id::text, user_id::text, day_of_week, start_time::text, end_time::text, timezone
		 FROM availability_slots WHERE user_id = ANY($1::uuid[])`,
		userIDs,
	)
}

func (c *PgSlotClient) querySlots(q string, arg any) ([]usecase.SlotData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := c.pool.Query(ctx, q, arg)
	if err != nil {
		return nil, fmt.Errorf("query slots: %w", err)
	}
	defer rows.Close()

	var out []usecase.SlotData
	for rows.Next() {
		var s usecase.SlotData
		var startStr, endStr string
		if err := rows.Scan(&s.ID, &s.UserID, &s.DayOfWeek, &startStr, &endStr, &s.Timezone); err != nil {
			return nil, err
		}
		// Trimming to HH:MM in case DB returns HH:MM:SS.
		s.StartTime = trimHHMM(startStr)
		s.EndTime = trimHHMM(endStr)
		out = append(out, s)
	}
	return out, rows.Err()
}

func trimHHMM(s string) string {
	if len(s) > 5 {
		return s[:5]
	}
	return s
}

// Course client
type PgCourseClient struct {
	pool *pgxpool.Pool
}

func NewPgCourseClient(pool *pgxpool.Pool) usecase.CourseClient {
	return &PgCourseClient{pool: pool}
}

func (c *PgCourseClient) ListCourseIDsForUser(userID string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	const q = `
SELECT id::text
FROM courses
WHERE owner_user_id = $1;
`
	rows, err := c.pool.Query(ctx, q, userID)
	if err != nil {
		return nil, fmt.Errorf("list course ids: %w", err)
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
