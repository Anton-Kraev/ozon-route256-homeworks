package order

import (
	"context"

	models "gitlab.ozon.dev/antonkraeww/homeworks/hw-3/internal/models/domain/order"
)

func (r OrderRepository) GetOrderByID(ctx context.Context, id uint64) (*models.Order, error) {
	const query = `SELECT * FROM orders WHERE id = $1`

	var order orderSchema

	row := r.pool.QueryRow(ctx, query, id)
	if err := row.Scan(&order); err != nil {
		return nil, err
	}

	ord := order.toDomain()

	return &ord, nil
}
