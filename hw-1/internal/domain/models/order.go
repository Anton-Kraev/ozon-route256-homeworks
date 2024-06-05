package models

import (
	"slices"
	"time"

	"gitlab.ozon.dev/antonkraeww/homeworks/hw-1/pkg/hash"
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

func (o *Order) MatchesFilter(filter OrderFilter) bool {
	matchesOrderID := len(filter.OrdersID) == 0 || slices.Contains(filter.OrdersID, o.OrderID)
	matchesClientID := len(filter.ClientsID) == 0 || slices.Contains(filter.ClientsID, o.ClientID)
	matchesStatus := len(filter.Statuses) == 0 || slices.Contains(filter.Statuses, o.Status)

	return matchesOrderID && matchesClientID && matchesStatus
}
