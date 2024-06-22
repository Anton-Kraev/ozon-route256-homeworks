package order

import (
	"time"
)

type Order struct {
	OrderID       uint64
	ClientID      uint64
	Weight        uint
	Cost          uint
	WrapType      Wrap
	StoredUntil   time.Time
	Status        Status
	StatusChanged time.Time
	Hash          string
}

func (o *Order) SetStatus(status Status, timeChanged time.Time) {
	o.Status = status
	o.StatusChanged = timeChanged
}

func (o *Order) SetHash(hash string) {
	o.Hash = hash
}
