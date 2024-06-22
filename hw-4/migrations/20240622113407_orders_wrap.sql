-- +goose Up
-- +goose StatementBegin
ALTER TABLE orders
    ADD COLUMN weight INTEGER CHECK (weight > 0) NOT NULL DEFAULT 1,
    ADD COLUMN cost INTEGER CHECK (cost > 0) NOT NULL DEFAULT 1,
    ADD COLUMN wrap_type TEXT NOT NULL DEFAULT 'Nowrap';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE orders
    DROP COLUMN weight,
    DROP COLUMN cost,
    DROP COLUMN wrap_type;
-- +goose StatementEnd
