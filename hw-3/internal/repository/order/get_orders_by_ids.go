package order

import (
	"context"
	"fmt"
	"strings"

	models "gitlab.ozon.dev/antonkraeww/homeworks/hw-3/internal/models/domain/order"
)

func (r OrderRepository) GetOrdersByIDs(ctx context.Context, ids []uint64) ([]models.Order, error) {
	placeholder := make([]string, len(ids))
	for i := range placeholder {
		placeholder[i] = fmt.Sprintf("$%d", i+1)
	}

	query := fmt.Sprintf("SELECT * FROM orders WHERE id IN (%s)", strings.Join(placeholder, ","))

	rows, err := r.pool.Query(ctx, query, ToInterfaceSlice(ids)...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		order  orderSchema
		orders []models.Order
	)

	for rows.Next() {
		if err = rows.Scan(&order); err != nil {
			return nil, err
		}

		orders = append(orders, order.toDomain())
	}

	return orders, nil
}

func ToInterfaceSlice[T any](slice []T) []interface{} {
	res := make([]interface{}, len(slice))
	for i, v := range slice {
		res[i] = v
	}

	return res
}
