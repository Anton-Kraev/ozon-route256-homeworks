package module

import (
	"time"

	errsdomain "gitlab.ozon.dev/antonkraeww/homeworks/hw-1/internal/domain/errors"
	"gitlab.ozon.dev/antonkraeww/homeworks/hw-1/internal/domain/models"
)

// ReturnOrder returns order to courier
func (m *OrderModule) ReturnOrder(orderID uint64) error {
	orders, err := m.Storage.GetOrders(models.OrderFilter{OrdersID: []uint64{orderID}})
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

	return m.Storage.ChangeOrders(map[uint64]models.Order{orderID: order})
}
