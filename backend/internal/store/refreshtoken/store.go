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
	ExpiresAt time.Time
	CreatedAt time.Time
}

type RefreshTokenStore interface {
	Create(context.Context, *RefreshToken) error
	GetByTokenHash(ctx context.Context, hash string) (*RefreshToken, error)
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
}

type PGRefreshTokenStore struct {
	pool *pgxpool.Pool
}

func NewPGRefreshTokenStore(pool *pgxpool.Pool) *PGRefreshTokenStore {
	return &PGRefreshTokenStore{pool: pool}
}

func (s *PGRefreshTokenStore) Create(ctx context.Context, token *RefreshToken) error {
	q := `INSERT INTO user_refresh_tokens 
			(user_id, token_hash, expires_at)
			VALUES 
			    ($1, $2, $3)`
	_, err := s.pool.Exec(ctx, q, token.UserID, token.TokenHash, token.ExpiresAt)
	return err
}

func (s *PGRefreshTokenStore) GetByTokenHash(ctx context.Context, hash string) (*RefreshToken, error) {
	q := `SELECT * FROM user_refresh_tokens 
         WHERE token_hash = $1 AND expires_at > now()`
	row := s.pool.QueryRow(ctx, q, hash)
	var t RefreshToken
	err := row.Scan(
		&t.ID,
		&t.UserID,
		&t.TokenHash,
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
