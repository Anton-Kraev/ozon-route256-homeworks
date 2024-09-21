package order

import (
	"context"

	"gitlab.ozon.dev/antonkraeww/ozon-route256-homeworks/internal/models/domain/order"
)

const timeFormat = "02.01.2006-15:04:05"

type orderService interface {
	RefundOrder(ctx context.Context, orderID, clientID uint64) error
	RefundsList(ctx context.Context, pageN, perPage uint) ([]order.Order, error)
	ReturnOrder(ctx context.Context, orderID uint64) error
	ClientOrders(ctx context.Context, clientID uint64, lastN uint, inStorage bool) ([]order.Order, error)
	ReceiveOrder(ctx context.Context, wrapType string, order order.Order) error
	DeliverOrders(ctx context.Context, ordersID []uint64) error
}

type OrderController struct {
	service orderService
}

func NewOrderController(service orderService) OrderController {
	return OrderController{service: service}
}
