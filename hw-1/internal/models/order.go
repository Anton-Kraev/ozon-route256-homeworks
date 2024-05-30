package models

import (
	"time"

	"gitlab.ozon.dev/antonkraeww/homeworks/hw-1/pkg/hash"
)

type Status string

const (
	Received  Status = "received"
	Returned  Status = "returned"
	Delivered Status = "delivered"
	Refunded  Status = "refunded"
)

type Order struct {
	OrderID       uint64
	ClientID      uint64
	StoredUntil   time.Time
	Status        Status
	StatusChanged time.Time
	Hash          string
}

func (o *Order) SetStatus(status Status, timeChanged time.Time) {
	o.Status = status
	o.StatusChanged = timeChanged
}

func (o *Order) SetHash() {
	o.Hash = hash.GenerateHash()
}
