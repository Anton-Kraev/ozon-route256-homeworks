package models

import "time"

type Status string

const (
	Received  Status = "received"
	Returned  Status = "returned"
	Delivered Status = "delivered"
	Refunded  Status = "refunded"
)

type Order struct {
	OrderID       int64
	ClientID      int64
	StoredUntil   time.Time
	Status        Status
	StatusChanged time.Time
	Hash          string
}
