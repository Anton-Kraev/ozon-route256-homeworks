package order

import (
	"context"

	models "gitlab.ozon.dev/antonkraeww/homeworks/hw-5/internal/models/domain/order"
	"gitlab.ozon.dev/antonkraeww/homeworks/hw-5/internal/pg"
)

func (r OrderRepository) ChangeOrders(ctx context.Context, changes []models.Order) error {
	tx, err := pg.GetTransactionFromContext(ctx)
	if err != nil {
		return err
	}

	const query = `
		UPDATE orders 
		SET status = $4, status_changed = $5, hash = $6 
		WHERE id = $1
	`

	for _, order := range changes {
		_, err = tx.Exec(
			ctx,
			query,
			order.Status,
			order.StatusChanged,
			order.Hash,
		)
		if err != nil {
			return err
		}
	}

	return nil
}
