package order

import (
	"time"

	"gitlab.ozon.dev/antonkraeww/homeworks/hw-5/internal/models/domain/wrap"
)

type Order struct {
	OrderID       uint64
	ClientID      uint64
	Weight        uint
	Cost          uint
	WrapType      string
	StoredUntil   time.Time
	Status        Status
	StatusChanged time.Time
	Hash          string
}

func (o *Order) Wrap(wrap wrap.Wrap) {
	o.WrapType = wrap.Name
	o.Weight += wrap.Weight
	o.Cost += wrap.Cost
}

func (o *Order) SetStatus(status Status, timeChanged time.Time) {
	o.Status = status
	o.StatusChanged = timeChanged
}

func (o *Order) SetHash(hash string) {
	o.Hash = hash
}
