-- +goose Up
-- +goose StatementBegin
UPDATE wrap SET max_weight = 29999 WHERE name = 'box';
UPDATE wrap SET max_weight = 9999 WHERE name = 'package';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
UPDATE wrap SET max_weight = 30 WHERE name = 'box';
UPDATE wrap SET max_weight = 10 WHERE name = 'package';
-- +goose StatementEnd
