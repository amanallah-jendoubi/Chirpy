-- +goose Up
-- +goose StatementBegin
ALTER TABLE messages
    ADD COLUMN body TEXT NOT NULL DEFAULT '';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE messages
    DROP COLUMN IF EXISTS body;
-- +goose StatementEnd