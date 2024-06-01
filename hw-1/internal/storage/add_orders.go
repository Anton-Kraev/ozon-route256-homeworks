package storage

import (
	errsdomain "gitlab.ozon.dev/antonkraeww/homeworks/hw-1/internal/domain/errors"
	"gitlab.ozon.dev/antonkraeww/homeworks/hw-1/internal/domain/models"
)

// AddOrders adds new orders to end of storage (if passed orders IDs is unique)
func (s OrderStorage) AddOrders(newOrders []models.Order) error {
	orders, err := s.readAll()
	if err != nil {
		return err
	}

	alreadyInStorage := make(map[uint64]struct{})
	for _, order := range orders {
		alreadyInStorage[order.OrderID] = struct{}{}
	}

	for _, order := range newOrders {
		if _, ok := alreadyInStorage[order.OrderID]; ok {
			return errsdomain.ErrOrderIDNotUnique(order.OrderID)
		}
	}

	return s.rewriteAll(append(orders, newOrders...))
}
