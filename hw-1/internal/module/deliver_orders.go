package module

import (
	"time"

	errsdomain "gitlab.ozon.dev/antonkraeww/homeworks/hw-1/internal/domain/errors"
	"gitlab.ozon.dev/antonkraeww/homeworks/hw-1/internal/domain/models"
)

// DeliverOrders deliver list of orders to client
func (m *OrderModule) DeliverOrders(ordersID []uint64) error {
	delivered := make(map[uint64]*models.Order)
	for _, orderID := range ordersID {
		delivered[orderID] = nil
	}

	orders, err := m.Storage.ReadAll()
	if err != nil {
		return err
	}

	var prevDelivered models.Order

	for _, order := range orders {
		if _, ok := delivered[order.OrderID]; !ok {
			continue
		}

		now := time.Now().UTC()

		if prevDelivered.OrderID != 0 && order.ClientID != prevDelivered.ClientID {
			return errsdomain.ErrDifferentClientOrders(order.ClientID, prevDelivered.ClientID)
		}
		if order.Status != models.Received {
			return errsdomain.ErrUnexpectedOrderStatus(order.OrderID, order.Status)
		}
		if now.After(order.StoredUntil) {
			return errsdomain.ErrRetentionPeriodExpired(order.OrderID)
		}

		order.SetStatus(models.Delivered, now)
		order.SetHash()

		prevDelivered = order
		delivered[order.OrderID] = &order
	}

	changes := make(map[uint64]models.Order)
	for id, order := range delivered {
		if order == nil {
			return errsdomain.ErrOrderNotFound(id)
		}
		changes[id] = *order
	}

	return m.Storage.ChangeOrders(changes)
}
