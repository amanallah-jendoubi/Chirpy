-- +goose Up
-- +goose StatementBegin
ALTER TABLE refresh_tokens
    ALTER COLUMN family_id SET DEFAULT gen_random_uuid();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE refresh_tokens
    ALTER COLUMN family_id DROP DEFAULT;
-- +goose StatementEnd