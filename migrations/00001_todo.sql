-- +goose Up
-- +goose StatementBegin
CREATE TABLE todo (
  id SERIAL PRIMARY KEY,
  text TEXT NOT NULL,
  done boolean
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
