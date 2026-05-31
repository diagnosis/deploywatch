package store

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DeviceToken struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	Token      string    `json:"token"`
	Platform   string    `json:"platform"`
	LastSeenAt time.Time `json:"last_seen_at"`
	CreatedAt  time.Time `json:"created_at"`
}

type DeviceTokenStore interface {
	Upsert(ctx context.Context, userID uuid.UUID, token, platform string) error
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]DeviceToken, error)
	Delete(ctx context.Context, token string) error
}

type PGDeviceTokenStore struct {
	pool *pgxpool.Pool
}

func NewPGDeviceTokenStore(pool *pgxpool.Pool) *PGDeviceTokenStore {
	return &PGDeviceTokenStore{pool: pool}
}

func (s *PGDeviceTokenStore) Upsert(ctx context.Context, userID uuid.UUID, token, platform string) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Remove this token from any other user first
	_, err = tx.Exec(ctx, `DELETE FROM device_tokens WHERE token = $1`, token)
	if err != nil {
		return err
	}

	// Remove old tokens for this user+platform
	_, err = tx.Exec(ctx, `DELETE FROM device_tokens WHERE user_id = $1 AND platform = $2`, userID, platform)
	if err != nil {
		return err
	}

	// Insert new token
	_, err = tx.Exec(ctx, `
        INSERT INTO device_tokens (user_id, token, platform, last_seen_at)
        VALUES ($1, $2, $3, now())
    `, userID, token, platform)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (s *PGDeviceTokenStore) GetByUserID(ctx context.Context, userID uuid.UUID) ([]DeviceToken, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, user_id, token, platform, last_seen_at, created_at
		FROM device_tokens WHERE user_id = $1
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokens []DeviceToken
	for rows.Next() {
		var t DeviceToken
		if err := rows.Scan(&t.ID, &t.UserID, &t.Token, &t.Platform, &t.LastSeenAt, &t.CreatedAt); err != nil {
			return nil, err
		}
		tokens = append(tokens, t)
	}
	return tokens, nil
}

func (s *PGDeviceTokenStore) Delete(ctx context.Context, token string) error {
	_, err := s.pool.Exec(ctx, `DELETE FROM device_tokens WHERE token = $1`, token)
	return err
}
