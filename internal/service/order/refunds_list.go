package order

import (
	"context"

	"github.com/Anton-Kraev/ozon-route256-homeworks/internal/models/domain/order"
)

// RefundsList returns list of refunds paginated
// optional pageN=<page number from the end>
// optional perPage=<number of orders per page>.
func (s *OrderService) RefundsList(ctx context.Context, pageN, perPage uint) ([]order.Order, error) {
	return s.orderRepo.GetOrdersByFilter(ctx, order.Filter{
		Statuses:     []order.Status{order.Refunded},
		PageN:        pageN,
		PerPage:      perPage,
		SortedByDate: true,
	})
}
