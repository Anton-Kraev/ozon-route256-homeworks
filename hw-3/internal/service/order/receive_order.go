package order

import (
	"context"
	"time"

	errsdomain "gitlab.ozon.dev/antonkraeww/homeworks/hw-3/internal/models/domain/errors"
	"gitlab.ozon.dev/antonkraeww/homeworks/hw-3/internal/models/domain/order"
)

// ReceiveOrder receives order from courier.
func (s *OrderService) ReceiveOrder(
	ctx context.Context, orderID, clientID uint64, storedUntil time.Time,
) error {
	now := time.Now().UTC()
	if now.After(storedUntil) {
		return errsdomain.ErrRetentionTimeInPast
	}

	newOrder := order.Order{
		OrderID:       orderID,
		ClientID:      clientID,
		StoredUntil:   storedUntil,
		Status:        order.Received,
		StatusChanged: now,
		Hash:          s.hashes.GetHash(),
	}

	return s.Repo.AddOrders([]order.Order{newOrder})
}
