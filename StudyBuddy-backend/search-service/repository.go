package main

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	pgvector "github.com/pgvector/pgvector-go"
)

// UserResult is the shape returned by SearchUsers.
// ID is a UUID string (matches users.id which is UUID).
type UserResult struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Bio       string `json:"bio"`
	Courses   string `json:"courses"`
	Interests string `json:"interests"`
}

// GetUserProfile reads the user's profile, interests, and courses from the DB
// and returns a single text string suitable for embedding.
// READ ONLY — never writes to the database.
func GetUserProfile(ctx context.Context, pool *pgxpool.Pool, userID string) (string, error) {
	// Each course is serialised as "Title (Subject, Level): Description" so the
	// embedding captures what the user actually studies, not just a bare title.
	const q = `
SELECT
    u.first_name,
    u.last_name,
    u.bio,
    COALESCE(string_agg(DISTINCT i.name, ', '), '') AS interests,
    COALESCE(
        string_agg(
            DISTINCT c.title || ' (' || c.subject || ', ' || c.level || '): ' || c.description,
            '; '
        ),
        ''
    ) AS courses
FROM users u
LEFT JOIN user_interests ui ON ui.user_id = u.id
LEFT JOIN interests i       ON i.id = ui.interest_id
LEFT JOIN courses c         ON c.owner_user_id = u.id
WHERE u.id = $1
GROUP BY u.id
`
	var firstName, lastName, bio, interests, courses string
	err := pool.QueryRow(ctx, q, userID).Scan(&firstName, &lastName, &bio, &interests, &courses)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("user %s not found", userID)
		}
		return "", fmt.Errorf("get user profile: %w", err)
	}

	var parts []string
	if name := strings.TrimSpace(firstName + " " + lastName); name != "" {
		parts = append(parts, "Name: "+name)
	}
	if bio != "" {
		parts = append(parts, "Bio: "+bio)
	}
	if interests != "" {
		parts = append(parts, "Interests: "+interests)
	}
	if courses != "" {
		parts = append(parts, "Courses: "+courses)
	}
	return strings.Join(parts, "\n"), nil
}

// UpsertEmbedding inserts or updates the embedding for a user.
func UpsertEmbedding(ctx context.Context, pool *pgxpool.Pool, userID string, embedding []float32) error {
	const q = `
INSERT INTO user_embeddings (user_id, embedding)
VALUES ($1, $2)
ON CONFLICT (user_id) DO UPDATE SET embedding = $2, updated_at = now()
`
	if _, err := pool.Exec(ctx, q, userID, pgvector.NewVector(embedding)); err != nil {
		return fmt.Errorf("upsert embedding: %w", err)
	}
	return nil
}

// SearchUsers finds the closest users to the given embedding vector using cosine distance.
func SearchUsers(ctx context.Context, pool *pgxpool.Pool, embedding []float32, limit int) ([]UserResult, error) {
	const q = `
SELECT
    u.id,
    u.first_name || ' ' || u.last_name                AS name,
    u.bio,
    COALESCE(string_agg(DISTINCT c.title, ', '), '')  AS courses,
    COALESCE(string_agg(DISTINCT i.name, ', '), '')   AS interests
FROM user_embeddings ue
JOIN users u            ON u.id = ue.user_id
LEFT JOIN user_interests ui ON ui.user_id = u.id
LEFT JOIN interests i       ON i.id = ui.interest_id
LEFT JOIN courses c         ON c.owner_user_id = u.id
WHERE u.is_active = true
GROUP BY u.id, ue.user_id
ORDER BY ue.embedding <=> $1
LIMIT $2
`
	rows, err := pool.Query(ctx, q, pgvector.NewVector(embedding), limit)
	if err != nil {
		return nil, fmt.Errorf("search users: %w", err)
	}
	defer rows.Close()

	results := make([]UserResult, 0)
	for rows.Next() {
		var r UserResult
		if err := rows.Scan(&r.ID, &r.Name, &r.Bio, &r.Courses, &r.Interests); err != nil {
			return nil, fmt.Errorf("scan result: %w", err)
		}
		results = append(results, r)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return results, nil
}
