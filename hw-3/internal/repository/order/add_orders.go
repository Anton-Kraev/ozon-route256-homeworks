package order

import (
	"fmt"

	errsdomain "gitlab.ozon.dev/antonkraeww/homeworks/hw-3/internal/models/domain/errors"
	models "gitlab.ozon.dev/antonkraeww/homeworks/hw-3/internal/models/domain/order"
)

// AddOrders adds new orders to end of storage (if passed orders IDs is unique).
func (r OrderRepository) AddOrders(newOrders []models.Order) error {
	orders, err := r.readAll()
	if err != nil {
		return err
	}

	alreadyInStorage := make(map[uint64]struct{})
	for _, order := range orders {
		alreadyInStorage[order.OrderID] = struct{}{}
	}

	for _, order := range newOrders {
		if _, ok := alreadyInStorage[order.OrderID]; ok {
			return fmt.Errorf("%w: %w",
				errsdomain.ErrOrderIDNotUnique,
				errsdomain.ErrorOrderIDNotUnique(order.OrderID),
			)
		}
	}

	return r.rewriteAll(append(orders, newOrders...))
}
