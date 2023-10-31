-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
ADD COLUMN google_id TEXT,
ALTER COLUMN password_hash DROP NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
