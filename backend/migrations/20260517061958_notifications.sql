-- +goose Up
-- +goose StatementBegin
CREATE TABLE notifications(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    read_at TIMESTAMPTZ,
    push_send_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT now(),
    CONSTRAINT uq_user_event UNIQUE  (user_id, event_id)
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS notifications;
-- +goose StatementEnd
