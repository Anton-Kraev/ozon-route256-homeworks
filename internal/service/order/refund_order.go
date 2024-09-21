package order

import (
	"context"
	"fmt"
	"time"

	errsdomain "github.com/Anton-Kraev/ozon-route256-homeworks/internal/models/domain/errors"
	models "github.com/Anton-Kraev/ozon-route256-homeworks/internal/models/domain/order"
)

// RefundOrder receives order refund from client.
func (s *OrderService) RefundOrder(ctx context.Context, orderID, clientID uint64) error {
	const maxRefundPeriod = time.Hour * 48

	order, err := s.orderRepo.GetOrderByID(ctx, orderID)
	if err != nil {
		return err
	}

	if order == nil || order.ClientID != clientID {
		return fmt.Errorf("%w: %w",
			errsdomain.ErrOrderNotFound,
			errsdomain.ErrorOrderNotFound(orderID),
		)
	}

	if order.Status == models.Refunded {
		return errsdomain.ErrOrderAlreadyRefunded
	}
	if order.Status != models.Delivered {
		return errsdomain.ErrOrderNotDeliveredYet
	}

	now := time.Now().UTC()
	if order.StatusChanged.Add(maxRefundPeriod).Before(now) {
		return errsdomain.ErrOrderDeliveredLongAgo
	}

	order.SetStatus(models.Refunded, now)
	order.SetHash(s.hashes.GetHash())

	return s.orderRepo.ChangeOrders(ctx, []models.Order{*order})
}
