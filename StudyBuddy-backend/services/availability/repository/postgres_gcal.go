package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"studybuddy/backend/pkg/crypto"
	"studybuddy/backend/services/availability/domain"
	"studybuddy/backend/services/availability/usecase"
)

type PgGCalRepository struct {
	pool *pgxpool.Pool
	key  []byte // 32-byte AES-256 encryption key
}

func NewPgGCalRepository(pool *pgxpool.Pool, encryptionKey []byte) usecase.GCalRepository {
	return &PgGCalRepository{pool: pool, key: encryptionKey}
}

func (r *PgGCalRepository) GetConnection(userID string) (*domain.GCalConnection, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	const q = `
SELECT user_id,
       access_token,
       refresh_token,
       token_expiry,
       calendar_id,
       sync_enabled,
       last_synced_at
FROM gcal_connections
WHERE user_id = $1;
`
	var (
		conn            domain.GCalConnection
		encAccessToken  string
		encRefreshToken string
	)

	err := r.pool.QueryRow(ctx, q, userID).Scan(
		&conn.UserID,
		&encAccessToken,
		&encRefreshToken,
		&conn.TokenExpiry,
		&conn.CalendarID,
		&conn.SyncEnabled,
		&conn.LastSyncedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("query gcal connection: %w", err)
	}

	conn.AccessToken, err = crypto.Decrypt(r.key, encAccessToken)
	if err != nil {
		return nil, fmt.Errorf("decrypt access_token: %w", err)
	}
	conn.RefreshToken, err = crypto.Decrypt(r.key, encRefreshToken)
	if err != nil {
		return nil, fmt.Errorf("decrypt refresh_token: %w", err)
	}

	return &conn, nil
}

func (r *PgGCalRepository) UpsertConnection(conn *domain.GCalConnection) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	encAccessToken, err := crypto.Encrypt(r.key, conn.AccessToken)
	if err != nil {
		return fmt.Errorf("encrypt access_token: %w", err)
	}
	encRefreshToken, err := crypto.Encrypt(r.key, conn.RefreshToken)
	if err != nil {
		return fmt.Errorf("encrypt refresh_token: %w", err)
	}

	const q = `
INSERT INTO gcal_connections
    (user_id, access_token, refresh_token, token_expiry, calendar_id, sync_enabled, last_synced_at)
VALUES
    ($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT (user_id) DO UPDATE
    SET access_token   = EXCLUDED.access_token,
        refresh_token  = EXCLUDED.refresh_token,
        token_expiry   = EXCLUDED.token_expiry,
        calendar_id    = EXCLUDED.calendar_id,
        sync_enabled   = EXCLUDED.sync_enabled,
        last_synced_at = EXCLUDED.last_synced_at,
        updated_at     = now();
`
	_, err = r.pool.Exec(ctx, q,
		conn.UserID,
		encAccessToken,
		encRefreshToken,
		conn.TokenExpiry,
		conn.CalendarID,
		conn.SyncEnabled,
		conn.LastSyncedAt, // may be nil
	)
	if err != nil {
		return fmt.Errorf("upsert gcal connection: %w", err)
	}
	return nil
}

func (r *PgGCalRepository) DeleteConnection(userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	const q = `DELETE FROM gcal_connections WHERE user_id = $1;`

	_, err := r.pool.Exec(ctx, q, userID)
	if err != nil {
		return fmt.Errorf("delete gcal connection: %w", err)
	}
	return nil
}
