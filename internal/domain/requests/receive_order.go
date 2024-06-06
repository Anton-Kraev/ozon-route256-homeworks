package requests

import (
	"time"
)

type ReceiveOrderRequest struct {
	OrderID     uint64
	ClientID    uint64
	StoredUntil time.Time
	Hash        string
}
