-- +goose Up
-- +goose StatementBegin
ALTER TABLE orders
    ADD COLUMN created_at TIMESTAMPTZ DEFAULT NOW(),
    ADD COLUMN updated_at TIMESTAMPTZ;

ALTER TABLE orders
    ALTER COLUMN stored_until TYPE TIMESTAMPTZ,
    ALTER COLUMN status_changed TYPE TIMESTAMPTZ;

ALTER TABLE orders
    RENAME COLUMN status_changed TO status_changed_at;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE orders
    DROP COLUMN updated_at,
    DROP COLUMN created_at;

ALTER TABLE orders
    ALTER COLUMN stored_until TYPE TIMESTAMP,
    ALTER COLUMN status_changed_at TYPE TIMESTAMP;

ALTER TABLE orders
    RENAME COLUMN status_changed_at TO status_changed;
-- +goose StatementEnd
