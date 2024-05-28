package storage

import (
	"time"

	"gitlab.ozon.dev/antonkraeww/homeworks/hw-1/internal/models"
)

type orderRecord struct {
	OrderID       int64     `json:"order_id"`
	ClientID      int64     `json:"client_id"`
	StoredUntil   time.Time `json:"stored_until"`
	Status        string    `json:"status"`
	StatusChanged time.Time `json:"status_changed"`
	Hash          string    `json:"hash"`
}

func (r orderRecord) toDomain() models.Order {
	return models.Order{
		OrderID:       r.OrderID,
		ClientID:      r.ClientID,
		StoredUntil:   r.StoredUntil,
		Status:        models.Status(r.Status),
		StatusChanged: r.StatusChanged,
		Hash:          r.Hash,
	}
}

// TODO: asynchronously add hash field
func toRecord(order models.Order) orderRecord {
	return orderRecord{
		OrderID:       order.OrderID,
		ClientID:      order.ClientID,
		StoredUntil:   order.StoredUntil,
		Status:        string(order.Status),
		StatusChanged: order.StatusChanged,
		Hash:          order.Hash,
	}
}
