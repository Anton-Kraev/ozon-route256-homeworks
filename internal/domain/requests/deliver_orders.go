package requests

type DeliverOrdersRequest struct {
	OrdersID []uint64
	Hashes   []string
}
