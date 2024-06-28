package order

import (
	"context"
	"fmt"
	"time"

	errsdomain "gitlab.ozon.dev/antonkraeww/homeworks/hw-5/internal/models/domain/errors"
	models "gitlab.ozon.dev/antonkraeww/homeworks/hw-5/internal/models/domain/order"
)

// ReturnOrder returns order to courier.
func (s *OrderService) ReturnOrder(ctx context.Context, orderID uint64) error {
	order, err := s.orderRepo.GetOrderByID(ctx, orderID)
	if err != nil {
		return err
	}

	if order == nil {
		return fmt.Errorf("%w: %w",
			errsdomain.ErrOrderNotFound,
			errsdomain.ErrorOrderNotFound(orderID),
		)
	}

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

	return s.orderRepo.ChangeOrders(ctx, []models.Order{*order})
}
