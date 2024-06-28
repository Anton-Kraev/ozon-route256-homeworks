package order

import (
	"context"

	"gitlab.ozon.dev/antonkraeww/homeworks/hw-5/internal/models/domain/order"
)

// ClientOrders returns list of client orders
// optional lastN for get last orders, by default return all orders
// optional inStorage for get only orders from storage.
func (s *OrderService) ClientOrders(
	ctx context.Context, clientID uint64, lastN uint, inStorage bool,
) ([]order.Order, error) {
	filter := order.Filter{
		ClientID:     clientID,
		PerPage:      lastN,
		SortedByDate: true,
	}
	if inStorage {
		filter.Statuses = []order.Status{order.Received, order.Refunded}
	}

	return s.orderRepo.GetOrdersByFilter(ctx, filter)
}
