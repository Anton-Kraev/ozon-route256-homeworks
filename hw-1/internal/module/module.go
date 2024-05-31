package module

import (
	"gitlab.ozon.dev/antonkraeww/homeworks/hw-1/internal/domain/models"
)

type orderStorage interface {
	AddOrder(newOrder models.Order) error
	ChangeOrders(changes map[uint64]models.Order) error
	FindOrder(orderID uint64) (*models.Order, error)
	ReadAll() ([]models.Order, error)
	RewriteAll(data []models.Order) error
}

type OrderModule struct {
	Storage orderStorage
}

func NewOrderModule(storage orderStorage) OrderModule {
	return OrderModule{Storage: storage}
}
