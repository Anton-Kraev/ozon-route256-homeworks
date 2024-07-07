-- +goose Up
-- +goose StatementBegin
ALTER TABLE wrap
    ADD COLUMN created_at TIMESTAMPTZ DEFAULT NOW(),
    ADD COLUMN updated_at TIMESTAMPTZ;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE wrap
    DROP COLUMN updated_at,
    DROP COLUMN created_at;
-- +goose StatementEnd
