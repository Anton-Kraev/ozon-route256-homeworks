package service

import (
	errsdomain "gitlab.ozon.dev/antonkraeww/homeworks/internal/models/domain/errors"
	"gitlab.ozon.dev/antonkraeww/homeworks/internal/models/domain/order"
	"gitlab.ozon.dev/antonkraeww/homeworks/internal/models/requests"
	"time"
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
		Hash:          req.Hash,
	}

	return s.Repo.AddOrders([]order.Order{newOrder})
}
