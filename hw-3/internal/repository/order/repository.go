package order

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

const table = "orders"

var columns = []string{"id", "client_id", "stored_until", "status", "status_changed", "hash"}

type OrderRepository struct {
	pool *pgxpool.Pool
}

func NewOrderRepository(pool *pgxpool.Pool) OrderRepository {
	return OrderRepository{pool: pool}
}
