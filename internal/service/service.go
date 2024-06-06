package service

import (
	"gitlab.ozon.dev/antonkraeww/homeworks/internal/domain/models"
)

type orderRepository interface {
	AddOrders(newOrders []models.Order) error
	ChangeOrders(changes map[uint64]models.Order) error
	GetOrders(filter models.OrderFilter) ([]models.Order, error)
}

type OrderService struct {
	Repo orderRepository
}

func NewOrderService(repo orderRepository) OrderService {
	return OrderService{Repo: repo}
}
