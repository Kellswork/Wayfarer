-- +goose Up
-- +goose StatementBegin
CREATE TABLE
  IF NOT EXISTS bus (
    id SERIAL PRIMARY KEY NOT NULL,
    plate_number VARCHAR(255) UNIQUE NOT NULL,
    manufacturer VARCHAR(255) NOT NULL,
    model VARCHAR(255) NOT NULL,
    type VARCHAR(255) NOT NULL,
    year VARCHAR(255) NOT NULL,
    capacity INT NOT NULL,

    created_at TIMESTAMP,
    updated_at TIMESTAMP
  );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS bus
-- +goose StatementEnd