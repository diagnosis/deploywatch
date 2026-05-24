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

type Event struct {
	ID         uuid.UUID `json:"id"`
	RepoID     int64     `json:"repo_id"`
	EventType  string    `json:"event_type"`
	Action     *string   `json:"action"`
	ActorLogin string    `json:"actor_login"`
	Payload    []byte    `json:"payload"`
	ReceivedAt time.Time `json:"received_at"`
}

type EventStore interface {
	Create(context.Context, *Event) (*Event, error)
	ListByRepoIDs(ctx context.Context, repoIDs []int64) ([]Event, error)
}

type PGEventStore struct {
	pool *pgxpool.Pool
}

func NewPGEventStore(pool *pgxpool.Pool) *PGEventStore {
	return &PGEventStore{pool}
}

func (s *PGEventStore) Create(ctx context.Context, e *Event) (*Event, error) {
	q := `INSERT INTO events (repo_id, event_type, action, actor_login, payload)
VALUES ($1, $2, $3, $4, $5)
RETURNING *`
	row := s.pool.QueryRow(ctx, q, e.RepoID, e.EventType, e.Action, e.ActorLogin, e.Payload)
	var e1 Event
	if err := row.Scan(
		&e1.ID,
		&e1.RepoID,
		&e1.EventType,
		&e1.Action,
		&e1.ActorLogin,
		&e1.Payload,
		&e1.ReceivedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors2.NotFound("user not found", "user not found")
		}
		return nil, err
	}
	return &e1, nil
}
func (s *PGEventStore) ListByRepoIDs(ctx context.Context, repoIDs []int64) ([]Event, error) {
	q := `SELECT * FROM events WHERE repo_id = ANY($1) ORDER BY received_at DESC LIMIT 50`
	rows, err := s.pool.Query(ctx, q, repoIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var e Event
		err := rows.Scan(&e.ID,
			&e.RepoID,
			&e.EventType,
			&e.Action,
			&e.ActorLogin,
			&e.Payload,
			&e.ReceivedAt)
		if err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, nil
}
