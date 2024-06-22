package order

import (
	"context"
	"time"

	errsdomain "gitlab.ozon.dev/antonkraeww/homeworks/hw-4/internal/models/domain/errors"
	models "gitlab.ozon.dev/antonkraeww/homeworks/hw-4/internal/models/domain/order"
	"gitlab.ozon.dev/antonkraeww/homeworks/hw-4/internal/models/domain/wrap"
)

// ReceiveOrder receives order from courier.
func (s *OrderService) ReceiveOrder(
	ctx context.Context, order models.Order,
) error {
	now := time.Now().UTC()
	if now.After(order.StoredUntil) {
		return errsdomain.ErrRetentionTimeInPast
	}

	wrapper, err := wrap.GetWrapper(order.WrapType)
	if err != nil {
		return err
	}

	if err = wrapper.WrapOrder(&order); err != nil {
		return err
	}

	orderWithSameID, err := s.Repo.GetOrderByID(ctx, order.OrderID)
	if err != nil {
		return err
	}
	if orderWithSameID != nil {
		return errsdomain.ErrOrderIDNotUnique
	}

	order.SetStatus(models.Received, now)
	order.SetHash(s.hashes.GetHash())

	return s.Repo.AddOrder(ctx, order)
}
