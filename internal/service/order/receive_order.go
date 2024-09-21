package order

import (
	"context"
	"time"

	errsdomain "github.com/Anton-Kraev/ozon-route256-homeworks/internal/models/domain/errors"
	models "github.com/Anton-Kraev/ozon-route256-homeworks/internal/models/domain/order"
)

// ReceiveOrder receives order from courier.
func (s *OrderService) ReceiveOrder(ctx context.Context, wrapType string, order models.Order) error {
	now := time.Now().UTC()
	if now.After(order.StoredUntil) {
		return errsdomain.ErrRetentionTimeInPast
	}

	orderWithSameID, err := s.orderRepo.GetOrderByID(ctx, order.OrderID)
	if err != nil {
		return err
	}
	if orderWithSameID != nil {
		return errsdomain.ErrOrderIDNotUnique
	}

	if wrapType != "" {
		wrap, err := s.wrapRepo.GetWrapByName(ctx, wrapType)
		if err != nil {
			return err
		}

		if wrap == nil {
			return errsdomain.ErrWrapNotFound
		}

		if err = order.Wrap(*wrap); err != nil {
			return err
		}
	}

	order.SetStatus(models.Received, now)
	order.SetHash(s.hashes.GetHash())

	return s.orderRepo.AddOrder(ctx, order)
}
