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

	const query = `
		INSERT INTO orders(id, client_id, weight, cost, wrap_type, stored_until, status, status_changed_at, hash) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err = tx.Exec(
		ctx,
		query,
		order.OrderID,
		order.ClientID,
		order.Weight,
		order.Cost,
		order.WrapType,
		order.StoredUntil,
		order.Status,
		order.StatusChanged,
		order.Hash,
	)

	return err
}
