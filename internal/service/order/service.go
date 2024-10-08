//go:generate mockgen -package=order -source=./service.go -destination=./service_mocks.go

package order

import (
	"context"

	"github.com/Anton-Kraev/ozon-route256-homeworks/internal/models/domain/order"
	"github.com/Anton-Kraev/ozon-route256-homeworks/internal/models/domain/wrap"
)

type (
	orderRepository interface {
		GetOrderByID(ctx context.Context, id uint64) (*order.Order, error)
		GetOrdersByIDs(ctx context.Context, ids []uint64) ([]order.Order, error)
		GetOrdersByFilter(ctx context.Context, filter order.Filter) ([]order.Order, error)
		AddOrder(ctx context.Context, order order.Order) error
		ChangeOrders(ctx context.Context, changes []order.Order) error
	}

	wrapRepository interface {
		GetWrapByName(ctx context.Context, name string) (*wrap.Wrap, error)
	}

	hashGenerator interface {
		GetHash() string
	}
)

type OrderService struct {
	orderRepo orderRepository
	wrapRepo  wrapRepository
	hashes    hashGenerator
}

func NewOrderService(orderRepo orderRepository, wrapRepo wrapRepository, hashGenerator hashGenerator) *OrderService {
	return &OrderService{
		orderRepo: orderRepo,
		wrapRepo:  wrapRepo,
		hashes:    hashGenerator,
	}
}
