-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS orders (
    id BIGSERIAL PRIMARY KEY,
    client_id INTEGER NOT NULL,
    stored_until TIMESTAMP NOT NULL,
    status TEXT NOT NULL,
    status_changed TIMESTAMP NOT NULL,
    hash VARCHAR(36) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders;
-- +goose StatementEnd
