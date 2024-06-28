package order

import (
	"database/sql"
	"time"

	"gitlab.ozon.dev/antonkraeww/homeworks/hw-5/internal/models/domain/order"
)

type orderSchema struct {
	OrderID         uint64         `db:"id"`
	ClientID        uint64         `db:"client_id"`
	Weight          uint           `db:"weight"`
	Cost            uint           `db:"cost"`
	WrapType        sql.NullString `db:"wrap_type"`
	StoredUntil     time.Time      `db:"stored_until"`
	Status          string         `db:"status"`
	StatusChangedAt time.Time      `db:"status_changed_at"`
	Hash            string         `db:"hash"`
	CreatedAt       time.Time      `db:"created_at"`
	UpdatedAt       sql.NullTime   `db:"updated_at"`
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
		StatusChanged: r.StatusChangedAt,
		Hash:          r.Hash,
	}
}
