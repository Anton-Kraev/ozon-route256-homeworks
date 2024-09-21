-- +goose Up
-- +goose StatementBegin
ALTER TABLE wrap ADD COLUMN max_weight INTEGER NOT NULL DEFAULT 2147483647;
UPDATE wrap SET max_weight = weight;
ALTER TABLE wrap
    DROP COLUMN weight,
    ADD CHECK (max_weight > 0);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE wrap ADD COLUMN weight INTEGER NOT NULL DEFAULT 1;
UPDATE wrap SET weight = max_weight;
ALTER TABLE wrap
    DROP COLUMN max_weight,
    ADD CHECK (weight > 0);
-- +goose StatementEnd
