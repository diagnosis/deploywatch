-- +goose Up
-- +goose StatementBegin
CREATE TABLE github_installations(
    installation_id BIGINT PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    account_login text NOT NULL,
    account_type text NOT NULL,
    suspended_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()

)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS github_installations;
-- +goose StatementEnd
