package order

import (
	"gitlab.ozon.dev/antonkraeww/homeworks/hw-2/internal/models/domain/order"
)

// RefundsList returns list of refunds paginated
// optional pageN=<page number from the end>
// optional perPage=<number of orders per page>.
func (s *OrderService) RefundsList(pageN, perPage uint) ([]order.Order, error) {
	orders, err := s.Repo.GetOrders(order.Filter{
		Statuses:     []order.Status{order.Refunded},
		PageN:        pageN,
		PerPage:      perPage,
		SortedByDate: true,
	})
	if err != nil {
		return []order.Order{}, err
	}

	return orders, nil
}
