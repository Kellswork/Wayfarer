-- +goose Up
-- +goose StatementBegin
CREATE TABLE
  IF NOT EXISTS trips (
    id uuid PRIMARY KEY NOT NULL,
    "busID" INT NOT NULL REFERENCES bus(id) ON DELETE CASCADE,
    origin VARCHAR(255) NOT NULL,
    destination VARCHAR(255) NOT NULL,
    trip_date TIMESTAMP NOT NULL,
    fare INT NOT NULL,
    status VARCHAR(255) NOT NULL,

    created_at TIMESTAMP,
    updated_at TIMESTAMP
  );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS trips;
-- +goose StatementEnd