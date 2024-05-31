package module

import (
	"time"

	domainErrors "gitlab.ozon.dev/antonkraeww/homeworks/hw-1/internal/domain/errors"
	"gitlab.ozon.dev/antonkraeww/homeworks/hw-1/internal/domain/models"
)

// ReturnOrder returns order to courier
func (m *OrderModule) ReturnOrder(orderID uint64) error {
	order, err := m.Storage.FindOrder(orderID)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	if order.StoredUntil.After(now) {
		return domainErrors.ErrRetentionPeriodNotExpiredYet
	}
	if order.Status == models.Returned {
		return domainErrors.ErrOrderAlreadyReturned
	}
	if order.Status == models.Delivered {
		return domainErrors.ErrOrderDelivered
	}

	order.SetStatus(models.Returned, now)
	order.SetHash()

	return m.Storage.ChangeOrders(map[uint64]models.Order{orderID: *order})
}
