package order

import (
	"context"

	models "github.com/Anton-Kraev/ozon-route256-homeworks/internal/models/domain/order"
	"github.com/Anton-Kraev/ozon-route256-homeworks/internal/pg"
)

func (r OrderRepository) ChangeOrders(ctx context.Context, changes []models.Order) error {
	tx, err := pg.GetTransactionFromContext(ctx)
	if err != nil {
		return err
	}

	const query = `
		UPDATE orders 
		SET status = $2, status_changed_at = $3, hash = $4, updated_at = $5
		WHERE id = $1
	`

	for _, order := range changes {
		_, err = tx.Exec(
			ctx,
			query,
			order.OrderID,
			order.Status,
			order.StatusChanged,
			order.Hash,
			order.StatusChanged,
		)
		if err != nil {
			return err
		}
	}

	return nil
}
