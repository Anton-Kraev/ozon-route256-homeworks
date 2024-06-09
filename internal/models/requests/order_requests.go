package requests

import "time"

type ClientOrdersRequest struct {
	ClientID  uint64
	LastN     uint
	InStorage bool
}

type DeliverOrdersRequest struct {
	OrdersID []uint64
}

type ReceiveOrderRequest struct {
	OrderID     uint64
	ClientID    uint64
	StoredUntil time.Time
}

type RefundOrderRequest struct {
	OrderID  uint64
	ClientID uint64
}

type RefundsListRequest struct {
	PageN   uint
	PerPage uint
}

type ReturnOrderRequest struct {
	OrderID uint64
}
