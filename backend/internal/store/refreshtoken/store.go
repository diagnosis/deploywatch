package store

import (
	"context"
	"errors"
	"time"

	errors2 "github.com/diagnosis/go-toolkit/errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RefreshToken struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	TokenHash string
	Platform  string
	ExpiresAt time.Time
	CreatedAt time.Time
}

type RefreshTokenStore interface {
	Create(ctx context.Context, t *RefreshToken) error
	GetByTokenHashAndPlatform(ctx context.Context, hash string, platform string) (*RefreshToken, error)
	GetByTokenHash(ctx context.Context, hash string) (*RefreshToken, error)
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
	DeleteByUserIDAndPlatform(ctx context.Context, userID uuid.UUID, platform string) error
}

type PGRefreshTokenStore struct {
	pool *pgxpool.Pool
}

func NewPGRefreshTokenStore(pool *pgxpool.Pool) *PGRefreshTokenStore {
	return &PGRefreshTokenStore{pool: pool}
}

func (s *PGRefreshTokenStore) Create(ctx context.Context, token *RefreshToken) error {
	q := `INSERT INTO user_refresh_tokens 
			(user_id, token_hash, platform, expires_at)
			VALUES 
			    ($1, $2, $3, $4)`
	_, err := s.pool.Exec(ctx, q, token.UserID, token.TokenHash, token.Platform, token.ExpiresAt)
	return err
}

func (s *PGRefreshTokenStore) GetByTokenHash(ctx context.Context, hash string) (*RefreshToken, error) {
	q := `SELECT id, user_id, token_hash, platform, expires_at, created_at
         FROM user_refresh_tokens 
         WHERE token_hash = $1 AND expires_at > now()`
	row := s.pool.QueryRow(ctx, q, hash)
	var t RefreshToken
	err := row.Scan(
		&t.ID,
		&t.UserID,
		&t.TokenHash,
		&t.Platform,
		&t.ExpiresAt,
		&t.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors2.NotFound("refresh token not found", "refresh token not found", err)
		}
		return nil, err
	}
	return &t, nil
}

func (s *PGRefreshTokenStore) GetByTokenHashAndPlatform(ctx context.Context, hash string, platform string) (*RefreshToken, error) {
	q := `SELECT id, user_id, token_hash, platform, expires_at, created_at
         FROM user_refresh_tokens 
         WHERE token_hash = $1 AND platform = $2 AND expires_at > now()`
	row := s.pool.QueryRow(ctx, q, hash, platform)
	var t RefreshToken
	err := row.Scan(
		&t.ID,
		&t.UserID,
		&t.TokenHash,
		&t.Platform,
		&t.ExpiresAt,
		&t.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors2.NotFound("refresh token not found", "refresh token not found", err)
		}
		return nil, err
	}
	return &t, nil
}

func (s *PGRefreshTokenStore) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	q := `DELETE FROM user_refresh_tokens WHERE user_id = $1`
	_, err := s.pool.Exec(ctx, q, userID)
	return err
}

func (s *PGRefreshTokenStore) DeleteByUserIDAndPlatform(ctx context.Context, userID uuid.UUID, platform string) error {
	q := `DELETE FROM user_refresh_tokens WHERE user_id = $1 AND platform = $2`
	_, err := s.pool.Exec(ctx, q, userID, platform)
	return err
}
