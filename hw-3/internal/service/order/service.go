package order

import (
	"context"

	"gitlab.ozon.dev/antonkraeww/homeworks/hw-3/internal/models/domain/order"
)

type orderRepository interface {
	GetOrderByID(ctx context.Context, id uint64) (*order.Order, error)
	GetOrdersByIDs(ctx context.Context, ids []uint64) ([]order.Order, error)
	GetOrdersByFilter(ctx context.Context, filter order.Filter) ([]order.Order, error)
	AddOrder(ctx context.Context, order order.Order) error
	ChangeOrders(ctx context.Context, changes []order.Order) error
}

type hashGenerator interface {
	GetHash() string
}

type OrderService struct {
	Repo   orderRepository
	hashes hashGenerator
}

func NewOrderService(repo orderRepository, hashGenerator hashGenerator) *OrderService {
	return &OrderService{Repo: repo, hashes: hashGenerator}
}
