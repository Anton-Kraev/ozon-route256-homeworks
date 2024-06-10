package service

import (
	"gitlab.ozon.dev/antonkraeww/homeworks/hw-2/internal/models/domain/order"
	"gitlab.ozon.dev/antonkraeww/homeworks/hw-2/internal/models/requests"
)

// RefundsList returns list of refunds paginated
// optional pageN=<page number from the end>
// optional perPage=<number of orders per page>.
func (s *OrderService) RefundsList(req requests.RefundsListRequest) ([]order.Order, error) {
	orders, err := s.Repo.GetOrders(order.Filter{
		Statuses:     []order.Status{order.Refunded},
		PageN:        req.PageN,
		PerPage:      req.PerPage,
		SortedByDate: true,
	})
	if err != nil {
		return []order.Order{}, err
	}

	return orders, nil
}
