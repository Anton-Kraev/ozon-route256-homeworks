package service

import (
	"time"

	errsdomain "gitlab.ozon.dev/antonkraeww/homeworks/internal/domain/errors"
	"gitlab.ozon.dev/antonkraeww/homeworks/internal/domain/models"
)

// ReceiveOrder receives order from courier.
func (s *OrderService) ReceiveOrder(orderID, clientID uint64, storedUntil time.Time) error {
	now := time.Now().UTC()
	if now.After(storedUntil) {
		return errsdomain.ErrRetentionTimeInPast
	}

	newOrder := models.Order{
		OrderID:       orderID,
		ClientID:      clientID,
		StoredUntil:   storedUntil,
		Status:        models.Received,
		StatusChanged: now,
	}
	newOrder.SetHash()

	return s.Repo.AddOrders([]models.Order{newOrder})
}
