package module

import (
	"time"

	errsdomain "gitlab.ozon.dev/antonkraeww/homeworks/hw-1/internal/domain/errors"
	"gitlab.ozon.dev/antonkraeww/homeworks/hw-1/internal/domain/models"
)

// RefundOrder receives order refund from client
func (m *OrderModule) RefundOrder(orderID, clientID uint64) error {
	orders, err := m.Storage.ReadAll()
	if err != nil {
		return err
	}

	for _, order := range orders {
		if order.OrderID != orderID || order.ClientID != clientID {
			continue
		}

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

		return m.Storage.ChangeOrders(map[uint64]models.Order{orderID: order})
	}

	return errsdomain.ErrOrderNotFound(orderID)
}
