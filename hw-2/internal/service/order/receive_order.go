package order

import (
	"time"

	errsdomain "gitlab.ozon.dev/antonkraeww/homeworks/hw-2/internal/models/domain/errors"
	"gitlab.ozon.dev/antonkraeww/homeworks/hw-2/internal/models/domain/order"
	"gitlab.ozon.dev/antonkraeww/homeworks/hw-2/internal/models/requests"
)

// ReceiveOrder receives order from courier.
func (s *OrderService) ReceiveOrder(req requests.ReceiveOrderRequest) error {
	now := time.Now().UTC()
	if now.After(req.StoredUntil) {
		return errsdomain.ErrRetentionTimeInPast
	}

	newOrder := order.Order{
		OrderID:       req.OrderID,
		ClientID:      req.ClientID,
		StoredUntil:   req.StoredUntil,
		Status:        order.Received,
		StatusChanged: now,
		Hash:          s.hashes.GetHash(),
	}

	return s.Repo.AddOrders([]order.Order{newOrder})
}
