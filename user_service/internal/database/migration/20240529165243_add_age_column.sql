-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
ADD COLUMN age date;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
DROP COLUMN age;
-- +goose StatementEnd
