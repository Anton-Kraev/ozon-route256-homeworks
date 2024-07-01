package order

import (
	"context"

	"gitlab.ozon.dev/antonkraeww/homeworks/hw-5/internal/models/domain/order"
	"gitlab.ozon.dev/antonkraeww/homeworks/hw-5/internal/pg"
)

func (r OrderRepository) AddOrder(ctx context.Context, order order.Order) error {
	tx, err := pg.GetTransactionFromContext(ctx)
	if err != nil {
		return err
	}

	const (
		insertQuery = `
			INSERT INTO orders(id, client_id, weight, cost, stored_until, status, status_changed_at, hash) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		`

		updateQuery = `UPDATE orders SET wrap_type = $1 WHERE id = $2`
	)

	_, err = tx.Exec(
		ctx,
		insertQuery,
		order.OrderID,
		order.ClientID,
		order.Weight,
		order.Cost,
		order.StoredUntil,
		order.Status,
		order.StatusChanged,
		order.Hash,
	)
	if err == nil && order.WrapType != "" {
		_, err = tx.Exec(ctx, updateQuery, order.WrapType, order.OrderID)
	}

	return err
}
