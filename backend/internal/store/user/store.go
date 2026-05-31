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

type User struct {
	ID          uuid.UUID `json:"id"`
	GitHubID    int64     `json:"git_hub_id"`
	Login       string    `json:"login"`
	Name        *string   `json:"name,omitempty"`
	AvatarURL   *string   `json:"avatar_url,omitempty"`
	Email       *string   `json:"email,omitempty"`
	AccessToken string    `json:"access_token"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UserStore interface {
	UpsertUser(context.Context, *User) (*User, error)
	GetUserByID(context.Context, uuid.UUID) (*User, error)
	GetUserByLogin(ctx context.Context, login string) (*User, error)
	DeleteUser(ctx context.Context, userID uuid.UUID) error
}

type PGUserStore struct {
	pool *pgxpool.Pool
}

func NewPGUserStore(pool *pgxpool.Pool) *PGUserStore { return &PGUserStore{pool} }

func (s *PGUserStore) UpsertUser(ctx context.Context, user *User) (*User, error) {
	q := `
		INSERT INTO users (github_id, login, name, avatar_url, email, access_token)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (github_id) 
		DO UPDATE SET
    	login = EXCLUDED.login,
    	name = EXCLUDED.name,
    	avatar_url = EXCLUDED.avatar_url,
    	email = EXCLUDED.email,
    	access_token = EXCLUDED.access_token,
    	updated_at = now()
		RETURNING *
`
	row := s.pool.QueryRow(
		ctx, q, user.GitHubID, user.Login, user.Name,
		user.AvatarURL, user.Email, user.AccessToken,
	)

	var u User
	err := row.Scan(
		&u.ID, &u.GitHubID, &u.Login, &u.Name, &u.AvatarURL,
		&u.Email, &u.AccessToken, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *PGUserStore) GetUserByID(ctx context.Context, id uuid.UUID) (*User, error) {
	q := `SELECT * FROM users WHERE id = $1`

	row := s.pool.QueryRow(ctx, q, id)

	var u User
	err := row.Scan(
		&u.ID, &u.GitHubID, &u.Login, &u.Name, &u.AvatarURL,
		&u.Email, &u.AccessToken, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors2.NotFound("user not found", "user not found")
		}
		return nil, err
	}
	return &u, nil

}
func (s *PGUserStore) GetUserByLogin(ctx context.Context, login string) (*User, error) {
	q := `SELECT * FROM users WHERE login = $1`

	row := s.pool.QueryRow(ctx, q, login)

	var u User
	err := row.Scan(
		&u.ID, &u.GitHubID, &u.Login, &u.Name, &u.AvatarURL,
		&u.Email, &u.AccessToken, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors2.NotFound("user not found", "user not found")
		}
		return nil, err
	}
	return &u, nil

}

func (s *PGUserStore) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	_, err := s.pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, userID)
	return err
}
