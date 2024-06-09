package service

import (
	"gitlab.ozon.dev/antonkraeww/homeworks/internal/models/domain/order"
)

type orderRepository interface {
	AddOrders(newOrders []order.Order) error
	ChangeOrders(changes map[uint64]order.Order) error
	GetOrders(filter order.OrderFilter) ([]order.Order, error)
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
