package service

import (
	"time"

	errsdomain "gitlab.ozon.dev/antonkraeww/homeworks/internal/domain/errors"
	"gitlab.ozon.dev/antonkraeww/homeworks/internal/domain/models"
	"gitlab.ozon.dev/antonkraeww/homeworks/internal/domain/requests"
)

// ReceiveOrder receives order from courier.
func (s *OrderService) ReceiveOrder(req requests.ReceiveOrderRequest) error {
	now := time.Now().UTC()
	if now.After(req.StoredUntil) {
		return errsdomain.ErrRetentionTimeInPast
	}

	newOrder := models.Order{
		OrderID:       req.OrderID,
		ClientID:      req.ClientID,
		StoredUntil:   req.StoredUntil,
		Status:        models.Received,
		StatusChanged: now,
		Hash:          req.Hash,
	}

	return s.Repo.AddOrders([]models.Order{newOrder})
}
