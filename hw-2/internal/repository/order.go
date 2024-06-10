package repository

import (
	"time"

	"gitlab.ozon.dev/antonkraeww/homeworks/hw-2/internal/models/domain/order"
)

type orderRecord struct {
	OrderID       uint64    `json:"order_id"`
	ClientID      uint64    `json:"client_id"`
	StoredUntil   time.Time `json:"stored_until"`
	Status        string    `json:"status"`
	StatusChanged time.Time `json:"status_changed"`
	Hash          string    `json:"hash"`
}

func (r orderRecord) toDomain() order.Order {
	return order.Order{
		OrderID:       r.OrderID,
		ClientID:      r.ClientID,
		StoredUntil:   r.StoredUntil,
		Status:        order.Status(r.Status),
		StatusChanged: r.StatusChanged,
		Hash:          r.Hash,
	}
}

func toRecord(order order.Order) orderRecord {
	return orderRecord{
		OrderID:       order.OrderID,
		ClientID:      order.ClientID,
		StoredUntil:   order.StoredUntil,
		Status:        string(order.Status),
		StatusChanged: order.StatusChanged,
		Hash:          order.Hash,
	}
}
