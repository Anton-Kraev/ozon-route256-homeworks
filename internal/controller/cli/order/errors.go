package order

import "errors"

var (
	errParseArgs     = errors.New("can't parse args")
	errClientOrders  = errors.New("can't get client orders")
	errDeliverOrders = errors.New("can't deliver orders")
	errReceiveOrder  = errors.New("can't receive order")
	errRefundOrder   = errors.New("can't refund order")
	errReturnOrder   = errors.New("can't return order")
	errRefundsList   = errors.New("can't refunds list")
)
