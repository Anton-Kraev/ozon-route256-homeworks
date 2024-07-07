-- +goose Up
-- +goose StatementBegin
CREATE TEMPORARY TABLE temp_orders AS
SELECT id, client_id, stored_until, status, status_changed_at, hash, weight, cost, created_at, updated_at, wrap_type
FROM orders;

ALTER TABLE orders
    DROP COLUMN id;
ALTER TABLE orders
    ADD COLUMN id BIGINT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY;

INSERT INTO
    orders(id, client_id, stored_until, status, status_changed_at, hash, weight, cost, created_at, updated_at, wrap_type)
    SELECT id, client_id, stored_until, status, status_changed_at, hash, weight, cost, created_at, updated_at, wrap_type
FROM temp_orders;

DROP TABLE temp_orders;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
CREATE TEMPORARY TABLE temp_orders AS
SELECT id, client_id, stored_until, status, status_changed_at, hash, weight, cost, created_at, updated_at, wrap_type
FROM orders;

ALTER TABLE orders
    DROP COLUMN id;

ALTER TABLE orders
    ADD COLUMN id BIGSERIAL PRIMARY KEY;

INSERT INTO
    orders(id, client_id, stored_until, status, status_changed_at, hash, weight, cost, created_at, updated_at, wrap_type)
    SELECT id, client_id, stored_until, status, status_changed_at, hash, weight, cost, created_at, updated_at, wrap_type
FROM temp_orders;

DROP TABLE temp_orders;
-- +goose StatementEnd
