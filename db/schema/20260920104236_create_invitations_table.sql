-- +goose Up
-- +goose StatementBegin
CREATE TABLE invitations (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    status            TEXT NOT NULL DEFAULT 'pending',
    conversation_id   UUID NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    sender_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    receiver_id       UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_invitations_conversation_id ON invitations(conversation_id);
CREATE INDEX idx_invitations_sender_id ON invitations(sender_id);
CREATE INDEX idx_invitations_receiver_id ON invitations(receiver_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS invitations;
-- +goose StatementEnd