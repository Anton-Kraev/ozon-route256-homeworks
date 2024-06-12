package order

import (
	"fmt"
	"time"

	errsdomain "gitlab.ozon.dev/antonkraeww/homeworks/hw-2/internal/models/domain/errors"
	models "gitlab.ozon.dev/antonkraeww/homeworks/hw-2/internal/models/domain/order"
)

// ReturnOrder returns order to courier.
func (s *OrderService) ReturnOrder(orderID uint64) error {
	orders, err := s.Repo.GetOrders(models.Filter{OrdersID: []uint64{orderID}})
	if err != nil {
		return err
	}

	if len(orders) == 0 {
		return fmt.Errorf("%w: %w",
			errsdomain.ErrOrderNotFound,
			errsdomain.ErrorOrderNotFound(orderID),
		)
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
	order.SetHash(s.hashes.GetHash())

	return s.Repo.ChangeOrders(map[uint64]models.Order{orderID: order})
}
