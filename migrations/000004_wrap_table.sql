-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS wrap (
    name TEXT PRIMARY KEY,
    weight INTEGER CHECK (weight > 0) NOT NULL,
    cost INTEGER CHECK (cost > 0) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS wrap;
-- +goose StatementEnd
