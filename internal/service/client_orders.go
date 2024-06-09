package service

import (
	"gitlab.ozon.dev/antonkraeww/homeworks/internal/models/domain/order"
	"gitlab.ozon.dev/antonkraeww/homeworks/internal/models/requests"
)

// ClientOrders returns list of client orders
// optional lastN for get last orders, by default return all orders
// optional inStorage for get only orders from storage.
func (s *OrderService) ClientOrders(req requests.ClientOrdersRequest) ([]order.Order, error) {
	filter := order.OrderFilter{
		ClientsID:    []uint64{req.ClientID},
		PerPage:      req.LastN,
		SortedByDate: true,
	}
	if req.InStorage {
		filter.Statuses = []order.Status{order.Received, order.Refunded}
	}

	return s.Repo.GetOrders(filter)
}
