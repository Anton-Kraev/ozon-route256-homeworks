-- +goose Up
-- +goose StatementBegin
INSERT INTO wrap(name, weight, cost)
VALUES ('box', 30, 20), ('package', 10, 5), ('tape', 2147483647, 1);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM wrap WHERE name IN ('box', 'package', 'tape');
-- +goose StatementEnd
