-- +goose Up
-- +goose StatementBegin
CREATE TABLE events(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    repo_id BIGINT NOT NULL,
    event_type TEXT NOT NULL,
    action TEXT,
    actor_login TEXT NOT NULL,
    payload JSONB NOT NULL,
    received_at TIMESTAMPTZ DEFAULT now()
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS events;
-- +goose StatementEnd
