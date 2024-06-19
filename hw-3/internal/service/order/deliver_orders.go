package order

import (
	"context"
	"fmt"
	"time"

	errsdomain "gitlab.ozon.dev/antonkraeww/homeworks/hw-3/internal/models/domain/errors"
	models "gitlab.ozon.dev/antonkraeww/homeworks/hw-3/internal/models/domain/order"
)

// DeliverOrders deliver list of orders to client.
func (s *OrderService) DeliverOrders(ctx context.Context, ordersID []uint64) error {
	orders, err := s.Repo.GetOrdersByIDs(ctx, ordersID)
	if err != nil {
		return err
	}

	if len(orders) != len(ordersID) {
		return errsdomain.ErrOrderNotFound
	}

	var (
		order     *models.Order
		prevOrder *models.Order
	)

	for i := 0; i < len(orders); i++ {
		order = &orders[i]

		now := time.Now().UTC()
		if now.After(order.StoredUntil) {
			return fmt.Errorf("%w: %w",
				errsdomain.ErrRetentionPeriodExpired,
				errsdomain.ErrorRetentionPeriodExpired(order.OrderID),
			)
		}

		if prevOrder != nil && order.ClientID != prevOrder.ClientID {
			return fmt.Errorf("%w: %w",
				errsdomain.ErrDifferentClients,
				errsdomain.ErrorDifferentClients(order.ClientID, prevOrder.ClientID),
			)
		}

		if order.Status != models.Received {
			return fmt.Errorf("%w: %w",
				errsdomain.ErrUnexpectedOrderStatus,
				errsdomain.ErrorUnexpectedOrderStatus(order.OrderID, order.Status),
			)
		}

		order.SetStatus(models.Delivered, now)
		order.SetHash(s.hashes.GetHash())

		prevOrder = order
	}

	return s.Repo.ChangeOrders(ctx, orders)
}
