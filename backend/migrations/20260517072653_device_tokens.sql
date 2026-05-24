-- +goose Up
-- +goose StatementBegin
CREATE TABLE device_tokens(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token TEXT UNIQUE NOT NULL,
    platform TEXT NOT NULL,
    last_seen_at TIMESTAMPTZ DEFAULT now(),
    created_at TIMESTAMPTZ DEFAULT now(),
    CONSTRAINT chk_platform CHECK (platform IN ('ios', 'android'))
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS device_tokens;
-- +goose StatementEnd
