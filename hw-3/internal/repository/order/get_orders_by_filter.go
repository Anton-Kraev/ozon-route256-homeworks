package order

import (
	"context"
	"fmt"
	"strings"

	"github.com/georgysavva/scany/v2/pgxscan"

	models "gitlab.ozon.dev/antonkraeww/homeworks/hw-3/internal/models/domain/order"
)

func (r OrderRepository) GetOrdersByFilter(ctx context.Context, filter models.Filter) ([]models.Order, error) {
	var (
		baseQuery                        = "SELECT * FROM orders"
		args                             []interface{}
		conditions                       []string
		wherePart, sortedPart, limitPart string
	)

	filter.Init()

	if filter.ClientID != 0 {
		conditions = append(conditions, "client_id = ?")
		args = append(args, filter.ClientID)
	}

	if filter.Statuses != nil {
		var placeholders []string

		for _, status := range filter.Statuses {
			placeholders = append(placeholders, "?")
			args = append(args, status)
		}

		conditions = append(conditions, fmt.Sprintf("status IN (%s)", strings.Join(placeholders, ",")))
	}

	if filter.SortedByDate {
		sortedPart = "ORDER BY status_changed DESC"
	}

	limitPart = "LIMIT ? OFFSET ?"
	args = append(args, filter.PerPage, (filter.PageN)*filter.PerPage)

	if len(conditions) > 0 {
		wherePart = "WHERE " + strings.Join(conditions, " AND ")
	}

	finalQuery := fmt.Sprintf("%s %s %s %s", baseQuery, wherePart, sortedPart, limitPart)

	for i := range args {
		finalQuery = strings.Replace(finalQuery, "?", fmt.Sprintf("$%d", i+1), 1)
	}

	rows, err := r.pool.Query(ctx, finalQuery, args...)
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
