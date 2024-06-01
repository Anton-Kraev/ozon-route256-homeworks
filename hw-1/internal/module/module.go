package module

import (
	"gitlab.ozon.dev/antonkraeww/homeworks/hw-1/internal/domain/models"
)

type orderStorage interface {
	AddOrders(newOrders []models.Order) error
	ChangeOrders(changes map[uint64]models.Order) error
	GetOrders(filter models.OrderFilter) ([]models.Order, error)
}

type OrderModule struct {
	Storage orderStorage
}

func NewOrderModule(storage orderStorage) OrderModule {
	return OrderModule{Storage: storage}
}
