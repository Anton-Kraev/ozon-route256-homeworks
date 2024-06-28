package order

import (
	"database/sql"
	"time"

	"gitlab.ozon.dev/antonkraeww/homeworks/hw-5/internal/models/domain/order"
)

type orderSchema struct {
	OrderID       uint64         `db:"id"`
	ClientID      uint64         `db:"client_id"`
	Weight        uint           `db:"weight"`
	Cost          uint           `db:"cost"`
	WrapType      sql.NullString `db:"wrap_type"`
	StoredUntil   time.Time      `db:"stored_until"`
	Status        string         `db:"status"`
	StatusChanged time.Time      `db:"status_changed"`
	Hash          string         `db:"hash"`
}

func (r orderSchema) toDomain() order.Order {
	return order.Order{
		OrderID:       r.OrderID,
		ClientID:      r.ClientID,
		Weight:        r.Weight,
		Cost:          r.Cost,
		WrapType:      r.WrapType.String,
		StoredUntil:   r.StoredUntil,
		Status:        order.Status(r.Status),
		StatusChanged: r.StatusChanged,
		Hash:          r.Hash,
	}
}
