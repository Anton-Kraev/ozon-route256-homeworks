package service

import (
	"time"

	errsdomain "gitlab.ozon.dev/antonkraeww/homeworks/internal/domain/errors"
	"gitlab.ozon.dev/antonkraeww/homeworks/internal/domain/models"
)

// RefundOrder receives order refund from client.
func (s *OrderService) RefundOrder(orderID, clientID uint64) error {
	orders, err := s.Repo.GetOrders(models.OrderFilter{
		OrdersID:  []uint64{orderID},
		ClientsID: []uint64{clientID},
	})
	if err != nil {
		return err
	}

	if len(orders) != 0 {
		return errsdomain.ErrOrderNotFound(orderID)
	}

	order := orders[0]
	if order.Status == models.Refunded {
		return errsdomain.ErrOrderAlreadyRefunded
	}
	if order.Status != models.Delivered {
		return errsdomain.ErrOrderNotDeliveredYet
	}

	now := time.Now().UTC()
	if order.StatusChanged.Add(time.Hour * 48).Before(now) {
		return errsdomain.ErrOrderDeliveredLongAgo
	}

	order.SetStatus(models.Refunded, now)
	order.SetHash()

	return s.Repo.ChangeOrders(map[uint64]models.Order{orderID: order})
}
