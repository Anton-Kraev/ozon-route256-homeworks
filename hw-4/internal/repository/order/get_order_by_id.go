package order

import (
	"context"

	"github.com/georgysavva/scany/v2/pgxscan"

	models "gitlab.ozon.dev/antonkraeww/homeworks/hw-4/internal/models/domain/order"
)

func (r OrderRepository) GetOrderByID(ctx context.Context, id uint64) (*models.Order, error) {
	const query = `SELECT * FROM orders WHERE id = $1`

	var orders []orderSchema

	err := pgxscan.Select(ctx, r.pool, &orders, query, id)
	if err != nil || len(orders) == 0 {
		return nil, err
	}

	ord := orders[0].toDomain()

	return &ord, nil
}
