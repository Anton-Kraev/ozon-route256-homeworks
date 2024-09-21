package order

import (
	"context"
	"fmt"
	"time"

	errsdomain "gitlab.ozon.dev/antonkraeww/ozon-route256-homeworks/internal/models/domain/errors"
	models "gitlab.ozon.dev/antonkraeww/ozon-route256-homeworks/internal/models/domain/order"
)

// DeliverOrders deliver list of orders to client.
func (s *OrderService) DeliverOrders(ctx context.Context, ordersID []uint64) error {
	orders, err := s.orderRepo.GetOrdersByIDs(ctx, ordersID)
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

	for i := range len(orders) {
		order = &orders[i]

		now := time.Now().UTC()
		if now.After(order.StoredUntil) {
			return fmt.Errorf("%w: %w",
				errsdomain.ErrRetentionPeriodExpired,
				errsdomain.ErrorRetentionPeriodExpired(order.OrderID),
			)
		}

		if order.Status != models.Received {
			return fmt.Errorf("%w: %w",
				errsdomain.ErrUnexpectedOrderStatus,
				errsdomain.ErrorUnexpectedOrderStatus(order.OrderID, string(order.Status)),
			)
		}

		if prevOrder != nil && order.ClientID != prevOrder.ClientID {
			return fmt.Errorf("%w: %w",
				errsdomain.ErrDifferentClients,
				errsdomain.ErrorDifferentClients(order.ClientID, prevOrder.ClientID),
			)
		}

		order.SetStatus(models.Delivered, now)
		order.SetHash(s.hashes.GetHash())

		prevOrder = order
	}

	return s.orderRepo.ChangeOrders(ctx, orders)
}
