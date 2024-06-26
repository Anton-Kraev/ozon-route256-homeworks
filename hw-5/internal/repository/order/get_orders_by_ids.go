package order

import (
	"context"
	"fmt"
	"strings"

	"github.com/georgysavva/scany/v2/pgxscan"

	"gitlab.ozon.dev/antonkraeww/homeworks/hw-5/internal/helpers"
	models "gitlab.ozon.dev/antonkraeww/homeworks/hw-5/internal/models/domain/order"
)

func (r OrderRepository) GetOrdersByIDs(ctx context.Context, ids []uint64) ([]models.Order, error) {
	placeholder := make([]string, len(ids))
	for i := range placeholder {
		placeholder[i] = fmt.Sprintf("$%d", i+1)
	}

	query := fmt.Sprintf("SELECT * FROM orders WHERE id IN (%s)", strings.Join(placeholder, ","))

	rows, err := r.pool.Query(ctx, query, helpers.TypedSliceToInterfaceSlice(ids)...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		orders     []models.Order
		rowScanner = pgxscan.NewRowScanner(rows)
	)

	for rows.Next() {
		var order orderSchema
		if err = rowScanner.Scan(&order); err != nil {
			return nil, err
		}

		orders = append(orders, order.toDomain())
	}

	return orders, nil
}
