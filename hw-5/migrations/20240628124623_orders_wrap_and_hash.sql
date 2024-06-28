-- +goose Up
-- +goose StatementBegin
ALTER TABLE orders
    ALTER COLUMN hash TYPE TEXT,
    DROP COLUMN wrap_type,
    ADD COLUMN wrap_type TEXT,
    ADD CONSTRAINT fk_wrap_type FOREIGN KEY (wrap_type) REFERENCES wrap(name) ON DELETE SET NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE orders
    ALTER COLUMN hash TYPE VARCHAR(36),
    DROP CONSTRAINT fk_wrap_type,
    DROP COLUMN wrap_type,
    ADD COLUMN wrap_type TEXT NOT NULL DEFAULT 'Nowrap';
-- +goose StatementEnd
