
-- +goose Up
-- +goose StatementBegin
CREATE TABLE members (
    user_id           UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    conversation_id   UUID NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, conversation_id)
);
 
CREATE INDEX idx_members_conversation_id ON members(conversation_id);
-- +goose StatementEnd
 
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS members;
-- +goose StatementEnd
