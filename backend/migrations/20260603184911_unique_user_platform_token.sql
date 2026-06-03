-- +goose Up
-- +goose StatementBegin
DELETE FROM user_refresh_tokens a USING user_refresh_tokens b
WHERE a.id > b.id AND a.user_id = b.user_id AND a.platform = b.platform;

ALTER TABLE user_refresh_tokens ADD CONSTRAINT unique_user_platform UNIQUE (user_id, platform);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE user_refresh_tokens DROP CONSTRAINT unique_user_platform;
-- +goose StatementEnd