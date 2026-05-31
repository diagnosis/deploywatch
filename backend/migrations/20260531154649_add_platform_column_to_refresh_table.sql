-- +goose Up
-- +goose StatementBegin
ALTER TABLE user_refresh_tokens ADD COLUMN platform TEXT NOT NULL DEFAULT 'web';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE user_refresh_tokens DROP COLUMN platform;
-- +goose StatementEnd
