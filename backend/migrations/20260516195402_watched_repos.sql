-- +goose Up
-- +goose StatementBegin
CREATE TABLE watched_repos(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    installation_id BIGINT NOT NULL REFERENCES github_installations(installation_id) ON DELETE CASCADE,
    repo_id BIGINT NOT NULL,
    repo_full_name text NOT NULL,
    event_filters JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    CONSTRAINT uq_user_repo UNIQUE (user_id, repo_id)
)

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS watched_repos;
-- +goose StatementEnd
