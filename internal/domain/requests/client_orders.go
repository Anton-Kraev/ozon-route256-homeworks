package requests

type ClientOrdersRequest struct {
	ClientID  uint64
	LastN     uint
	InStorage bool
}
