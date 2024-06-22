package order

import (
	"context"
	"time"

	errsdomain "gitlab.ozon.dev/antonkraeww/homeworks/hw-4/internal/models/domain/errors"
	models "gitlab.ozon.dev/antonkraeww/homeworks/hw-4/internal/models/domain/order"
)

// ReceiveOrder receives order from courier.
func (s *OrderService) ReceiveOrder(
	ctx context.Context, orderID, clientID uint64, storedUntil time.Time,
) error {
	now := time.Now().UTC()
	if now.After(storedUntil) {
		return errsdomain.ErrRetentionTimeInPast
	}

	order, err := s.Repo.GetOrderByID(ctx, orderID)
	if err != nil {
		return err
	}
	if order == nil {
		return errsdomain.ErrOrderIDNotUnique
	}

	newOrder := models.Order{
		OrderID:       orderID,
		ClientID:      clientID,
		StoredUntil:   storedUntil,
		Status:        models.Received,
		StatusChanged: now,
		Hash:          s.hashes.GetHash(),
	}

	return s.Repo.AddOrder(ctx, newOrder)
}
