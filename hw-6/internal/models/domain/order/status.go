package order

type Status string

const (
	Received  Status = "received"
	Returned  Status = "returned"
	Delivered Status = "delivered"
	Refunded  Status = "refunded"
)
