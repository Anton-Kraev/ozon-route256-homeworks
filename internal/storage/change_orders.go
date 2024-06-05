package storage

import (
	errsdomain "gitlab.ozon.dev/antonkraeww/homeworks/internal/domain/errors"
	"gitlab.ozon.dev/antonkraeww/homeworks/internal/domain/models"
)

// ChangeOrders changes orders data in storage, key=<order id to change> value=<new order data>.
func (s OrderStorage) ChangeOrders(changes map[uint64]models.Order) error {
	orders, err := s.readAll()
	if err != nil {
		return err
	}

	for i, order := range orders {
		if _, ok := changes[order.OrderID]; !ok {
			continue
		}

		orders[i] = changes[order.OrderID]
		delete(changes, order.OrderID)
	}

	// return err with first order that not found
	if len(changes) != 0 {
		for orderID := range changes {
			return errsdomain.ErrOrderNotFound(orderID)
		}
	}

	return s.rewriteAll(orders)
}
