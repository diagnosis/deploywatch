package installation

import (
	"context"
	"errors"
	"time"

	errors2 "github.com/diagnosis/go-toolkit/errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Installation struct {
	InstallationID int64      `json:"installation_id"`
	UserID         uuid.UUID  `json:"user_id"`
	AccountLogin   string     `json:"account_login"`
	AccountType    string     `json:"account_type"`
	SuspendedAt    *time.Time `json:"suspended_at"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type InstallationStore interface {
	Upsert(context.Context, *Installation) (*Installation, error)
	Delete(ctx context.Context, installationID int64) error
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*Installation, error)
	Suspend(ctx context.Context, installationID int64) error
	Unsuspend(ctx context.Context, installationID int64) error
}

type PGInstallationStore struct {
	pool *pgxpool.Pool
}

func NewPGInstallationStore(pool *pgxpool.Pool) *PGInstallationStore {
	return &PGInstallationStore{pool: pool}
}

func (s *PGInstallationStore) Upsert(ctx context.Context, inst *Installation) (*Installation, error) {
	q := `INSERT INTO github_installations (installation_id, user_id, account_login, account_type)
VALUES ($1, $2, $3, $4)
ON CONFLICT (installation_id)
DO UPDATE SET
    account_login = EXCLUDED.account_login,
    account_type = EXCLUDED.account_type,
    updated_at = now()
RETURNING *`
	row := s.pool.QueryRow(ctx, q, inst.InstallationID, inst.UserID, inst.AccountLogin, inst.AccountType)
	var i Installation
	if err := row.Scan(
		&i.InstallationID,
		&i.UserID,
		&i.AccountLogin,
		&i.AccountType,
		&i.SuspendedAt,
		&i.CreatedAt,
		&i.UpdatedAt); err != nil {
		if errors.Is(pgx.ErrNoRows, err) {
			return nil, errors2.NotFound("row not found", "row not found")
		}
		return nil, err
	}
	return &i, nil
}

func (s *PGInstallationStore) Delete(ctx context.Context, instID int64) error {
	q := `DELETE from github_installations where installation_id = $1`
	ct, err := s.pool.Exec(ctx, q, instID)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return errors2.NotFound("row not found", "row not found")
	}
	return nil
}

func (s *PGInstallationStore) GetByUserID(ctx context.Context, userId uuid.UUID) ([]*Installation, error) {
	q := `SELECT * FROM github_installations WHERE user_id = $1`
	var insts []*Installation
	rows, err := s.pool.Query(ctx, q, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var i Installation
		if err := rows.Scan(
			&i.InstallationID,
			&i.UserID,
			&i.AccountLogin,
			&i.AccountType,
			&i.SuspendedAt,
			&i.CreatedAt,
			&i.UpdatedAt); err != nil {
			if errors.Is(pgx.ErrNoRows, err) {
				return nil, errors2.NotFound("row not found", "row not found")
			}
			return nil, err
		}
		insts = append(insts, &i)
	}
	return insts, nil
}

func (s *PGInstallationStore) Suspend(ctx context.Context, installationID int64) error {
	q := `UPDATE github_installations SET suspended_at = now() where installation_id = $1`
	ct, err := s.pool.Exec(ctx, q, installationID)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return errors2.NotFound("row not found", "row not found")
	}
	return nil
}
func (s *PGInstallationStore) Unsuspend(ctx context.Context, installationID int64) error {
	q := `UPDATE github_installations SET suspended_at = NULL where installation_id = $1`
	ct, err := s.pool.Exec(ctx, q, installationID)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return errors2.NotFound("row not found", "row not found")
	}
	return nil
}
