package requests

type RefundOrderRequest struct {
	OrderID  uint64
	ClientID uint64
	Hash     string
}
