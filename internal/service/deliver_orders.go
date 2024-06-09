package service

import (
	errsdomain "gitlab.ozon.dev/antonkraeww/homeworks/internal/models/domain/errors"
	models "gitlab.ozon.dev/antonkraeww/homeworks/internal/models/domain/order"
	"gitlab.ozon.dev/antonkraeww/homeworks/internal/models/requests"
	"time"
)

// DeliverOrders deliver list of orders to client.
func (s *OrderService) DeliverOrders(req requests.DeliverOrdersRequest) error {
	delivered := make(map[uint64]models.Order)
	for _, orderID := range req.OrdersID {
		delivered[orderID] = models.Order{}
	}

	orders, err := s.Repo.GetOrders(models.OrderFilter{OrdersID: req.OrdersID})
	if err != nil {
		return err
	}

	var prevDelivered models.Order

	for i, order := range orders {
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
		order.SetHash(req.Hashes[i])

		prevDelivered = order
		delivered[order.OrderID] = order
	}

	changes := make(map[uint64]models.Order)
	for id, order := range delivered {
		if order.OrderID == 0 {
			return errsdomain.ErrOrderNotFound(id)
		}
		changes[id] = order
	}

	return s.Repo.ChangeOrders(changes)
}
