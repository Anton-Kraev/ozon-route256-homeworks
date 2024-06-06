package service

import (
	"time"

	errsdomain "gitlab.ozon.dev/antonkraeww/homeworks/internal/domain/errors"
	"gitlab.ozon.dev/antonkraeww/homeworks/internal/domain/models"
)

// ReturnOrder returns order to courier.
func (s *OrderService) ReturnOrder(orderID uint64) error {
	orders, err := s.Repo.GetOrders(models.OrderFilter{OrdersID: []uint64{orderID}})
	if err != nil {
		return err
	}

	if len(orders) == 0 {
		return errsdomain.ErrOrderNotFound(orderID)
	}

	order := orders[0]
	now := time.Now().UTC()

	if order.StoredUntil.After(now) {
		return errsdomain.ErrRetentionPeriodNotExpiredYet
	}
	if order.Status == models.Returned {
		return errsdomain.ErrOrderAlreadyReturned
	}
	if order.Status == models.Delivered {
		return errsdomain.ErrOrderDelivered
	}

	order.SetStatus(models.Returned, now)
	order.SetHash()

	return s.Repo.ChangeOrders(map[uint64]models.Order{orderID: order})
}
