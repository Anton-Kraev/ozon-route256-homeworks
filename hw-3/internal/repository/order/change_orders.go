package order

import (
	"fmt"

	errsdomain "gitlab.ozon.dev/antonkraeww/homeworks/hw-3/internal/models/domain/errors"
	models "gitlab.ozon.dev/antonkraeww/homeworks/hw-3/internal/models/domain/order"
)

// ChangeOrders changes orders data in storage, key=<order id to change> value=<new order data>.
func (r OrderRepository) ChangeOrders(changes map[uint64]models.Order) error {
	orders, err := r.readAll()
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
			return fmt.Errorf("%w: %w",
				errsdomain.ErrOrderNotFound,
				errsdomain.ErrorOrderNotFound(orderID),
			)
		}
	}

	return r.rewriteAll(orders)
}
