package repository

import (
	"context"
	"errors"
	"strconv"
	"time"

	"studybuddy/backend/services/courses/domain"
	"studybuddy/backend/services/courses/usecase"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PgCourseRepository implements CourseRepository using PostgreSQL.
type PgCourseRepository struct {
	pool *pgxpool.Pool
}

// NewPgCourseRepository creates a new PgCourseRepository.
func NewPgCourseRepository(pool *pgxpool.Pool) usecase.CourseRepository {
	return &PgCourseRepository{pool: pool}
}

func (r *PgCourseRepository) Create(course *domain.Course) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	const q = `
INSERT INTO courses (title, description, subject, level, owner_user_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, created_at, updated_at;
`
	err := r.pool.QueryRow(ctx, q,
		course.Title,
		course.Description,
		course.Subject,
		course.Level,
		course.OwnerUserID,
	).Scan(&course.ID, &course.CreatedAt, &course.UpdatedAt)
	return err
}

func (r *PgCourseRepository) Update(course *domain.Course) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	const q = `
UPDATE courses
SET title = $2,
    description = $3,
    subject = $4,
    level = $5,
    updated_at = now()
WHERE id = $1;
`
	ct, err := r.pool.Exec(ctx, q,
		course.ID,
		course.Title,
		course.Description,
		course.Subject,
		course.Level,
	)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *PgCourseRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	const q = `
DELETE FROM courses
WHERE id = $1;
`
	_, err := r.pool.Exec(ctx, q, id)
	return err
}

func (r *PgCourseRepository) GetByID(id string) (*domain.Course, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	const q = `
SELECT id, title, description, subject, level, owner_user_id, created_at, updated_at
FROM courses
WHERE id = $1;
`
	var c domain.Course
	err := r.pool.QueryRow(ctx, q, id).Scan(
		&c.ID,
		&c.Title,
		&c.Description,
		&c.Subject,
		&c.Level,
		&c.OwnerUserID,
		&c.CreatedAt,
		&c.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

func (r *PgCourseRepository) List(filter usecase.ListCoursesFilter) ([]domain.Course, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Simple filtering on subject and level; pagination via limit/offset.
	const base = `
SELECT id, title, description, subject, level, owner_user_id, created_at, updated_at
FROM courses
WHERE 1=1
`

	// Build query dynamically but safely.
	args := []any{}
	where := ""
	argPos := 1
	if filter.Subject != "" {
		where += " AND subject = $" + strconv.Itoa(argPos)
		args = append(args, filter.Subject)
		argPos++
	}
	if filter.Level != "" {
		where += " AND level = $" + strconv.Itoa(argPos)
		args = append(args, filter.Level)
		argPos++
	}
	limit := filter.Limit
	if limit <= 0 {
		limit = 20
	}
	offset := filter.Offset
	if offset < 0 {
		offset = 0
	}
	where += " ORDER BY created_at DESC"
	where += " LIMIT $" + strconv.Itoa(argPos)
	args = append(args, limit)
	argPos++
	where += " OFFSET $" + strconv.Itoa(argPos)
	args = append(args, offset)

	rows, err := r.pool.Query(ctx, base+where, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.Course
	for rows.Next() {
		var c domain.Course
		if err := rows.Scan(
			&c.ID,
			&c.Title,
			&c.Description,
			&c.Subject,
			&c.Level,
			&c.OwnerUserID,
			&c.CreatedAt,
			&c.UpdatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return out, nil
}
