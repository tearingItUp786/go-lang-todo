-- +goose Up
-- +goose StatementBegin
CREATE TABLE sessions (
  id SERIAL PRIMARY KEY,
  user_id INT UNIQUE REFERENCES users (id) ON DELETE CASCADE,
  token_hash TEXT UNIQUE NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE sessions;
-- +goose StatementEnd
