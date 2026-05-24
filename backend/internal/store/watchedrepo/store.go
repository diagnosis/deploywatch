package watchedrepo

import (
	"context"
	"errors"
	"time"

	errors2 "github.com/diagnosis/go-toolkit/errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WatchedRepo struct {
	ID             uuid.UUID `json:"id"`
	UserID         uuid.UUID `json:"user_id"`
	InstallationID int64     `json:"installation_id"`
	RepoID         int64     `json:"repo_id"`
	RepoFullName   string    `json:"repo_full_name"`
	EventFilters   []byte    `json:"event_filters"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type WatchedRepoStore interface {
	Add(ctx context.Context, repo *WatchedRepo) (*WatchedRepo, error)
	Remove(ctx context.Context, userID uuid.UUID, repoId int64) error
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*WatchedRepo, error)
	GetUsersByRepoID(ctx context.Context, repoID int64) ([]uuid.UUID, error)
}

type PGWatchedRepoStore struct {
	pool *pgxpool.Pool
}

func NewPGWatchedRepoStore(pool *pgxpool.Pool) *PGWatchedRepoStore {
	return &PGWatchedRepoStore{pool: pool}
}

func (s *PGWatchedRepoStore) Add(ctx context.Context, repo *WatchedRepo) (*WatchedRepo, error) {
	q := `INSERT into watched_repos (user_id, installation_id, repo_id, repo_full_name, event_filters)
			VALUES ($1,$2,$3,$4,$5) 
			RETURNING *`
	row := s.pool.QueryRow(ctx, q, repo.UserID, repo.InstallationID, repo.RepoID, repo.RepoFullName, repo.EventFilters)
	var r WatchedRepo
	if err := row.Scan(
		&r.ID,
		&r.UserID,
		&r.InstallationID,
		&r.RepoID,
		&r.RepoFullName,
		&r.EventFilters,
		&r.CreatedAt,
		&r.UpdatedAt,
	); err != nil {
		if errors.Is(pgx.ErrNoRows, err) {
			return nil, errors2.NotFound("repo not found", "repo not found")
		}
		return nil, err
	}
	return &r, nil
}

func (s *PGWatchedRepoStore) Remove(ctx context.Context, userID uuid.UUID, repoID int64) error {
	q := "DELETE FROM watched_repos WHERE user_id = $1 AND repo_id = $2"
	ct, err := s.pool.Exec(ctx, q, userID, repoID)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return errors2.NotFound("repo not found", "repo not found")
	}
	return nil
}

func (s *PGWatchedRepoStore) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*WatchedRepo, error) {
	q := "SELECT * FROM watched_repos where user_id = $1"
	rows, err := s.pool.Query(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var repos []*WatchedRepo
	for rows.Next() {
		var r WatchedRepo
		if err := rows.Scan(
			&r.ID,
			&r.UserID,
			&r.InstallationID,
			&r.RepoID,
			&r.RepoFullName,
			&r.EventFilters,
			&r.CreatedAt,
			&r.UpdatedAt,
		); err != nil {
			if errors.Is(pgx.ErrNoRows, err) {
				return nil, errors2.NotFound("repo not found", "repo not found")
			}
			return nil, err
		}
		repos = append(repos, &r)
	}
	return repos, nil
}

func (s *PGWatchedRepoStore) GetUsersByRepoID(ctx context.Context, repoID int64) ([]uuid.UUID, error) {
	q := `SELECT user_id FROM watched_repos WHERE repo_id = $1`
	rows, err := s.pool.Query(ctx, q, repoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userIDs []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		userIDs = append(userIDs, id)
	}
	return userIDs, nil
}
