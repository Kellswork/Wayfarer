-- +goose Up
-- +goose StatementBegin
ALTER TABLE users 
  ADD CONSTRAINT unique_email UNIQUE (email);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP CONSTRAINT IF EXISTS unique_email;
-- +goose StatementEnd
