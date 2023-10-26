-- +goose Up
-- +goose StatementBegin
ALTER TABLE todo RENAME TO todos;
ALTER TABLE todos
ADD COLUMN user_id INTEGER NOT NULL,
ADD CONSTRAINT fk_user_id
FOREIGN KEY (user_id)
REFERENCES users(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE todos;
-- +goose StatementEnd
