package order

import (
	"gitlab.ozon.dev/antonkraeww/homeworks/hw-2/internal/models/domain/order"
)

type orderRepository interface {
	AddOrders(newOrders []order.Order) error
	ChangeOrders(changes map[uint64]order.Order) error
	GetOrders(filter order.Filter) ([]order.Order, error)
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
