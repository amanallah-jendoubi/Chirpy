-- +goose Up
-- +goose StatementBegin
CREATE TABLE chat_group_members (
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    chat_group_id   UUID NOT NULL REFERENCES chat_groups(id) ON DELETE CASCADE,
    joined_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, chat_group_id)
);

CREATE INDEX idx_chat_group_members_chat_group_id ON chat_group_members (chat_group_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS chat_group_members;
-- +goose StatementEnd