package order

import (
	"fmt"
	"time"

	errsdomain "gitlab.ozon.dev/antonkraeww/homeworks/hw-3/internal/models/domain/errors"
	models "gitlab.ozon.dev/antonkraeww/homeworks/hw-3/internal/models/domain/order"
)

// DeliverOrders deliver list of orders to client.
func (s *OrderService) DeliverOrders(ordersID []uint64) error {
	toDeliver := make(map[uint64]models.Order)
	for _, orderID := range ordersID {
		toDeliver[orderID] = models.Order{}
	}

	orders, err := s.Repo.GetOrders(models.Filter{OrdersID: ordersID})
	if err != nil {
		return err
	}

	var prevOrder models.Order

	for _, order := range orders {
		now := time.Now().UTC()

		if prevOrder.OrderID != 0 && order.ClientID != prevOrder.ClientID {
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
		if now.After(order.StoredUntil) {
			return fmt.Errorf("%w: %w",
				errsdomain.ErrRetentionPeriodExpired,
				errsdomain.ErrorRetentionPeriodExpired(order.OrderID),
			)
		}

		order.SetStatus(models.Delivered, now)
		order.SetHash(s.hashes.GetHash())

		prevOrder = order
		toDeliver[order.OrderID] = order
	}

	return s.Repo.ChangeOrders(toDeliver)
}
