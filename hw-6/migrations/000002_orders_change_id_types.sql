-- +goose Up
-- +goose StatementBegin
ALTER TABLE orders
    ALTER COLUMN id TYPE BIGINT,
    ALTER COLUMN client_id TYPE BIGINT,
    ADD CONSTRAINT id_positive CHECK (id > 0),
    ADD CONSTRAINT client_id_positive CHECK (client_id > 0);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE orders
    DROP CONSTRAINT IF EXISTS id_positive,
    DROP CONSTRAINT IF EXISTS client_id_positive;

ALTER TABLE orders
    ALTER COLUMN id TYPE BIGSERIAL,
    ALTER COLUMN client_id TYPE INTEGER;

CREATE SEQUENCE IF NOT EXISTS orders_id_seq;
ALTER TABLE orders ALTER COLUMN id SET DEFAULT nextval('orders_id_seq');
ALTER SEQUENCE orders_id_seq OWNED BY orders.id;
-- +goose StatementEnd
