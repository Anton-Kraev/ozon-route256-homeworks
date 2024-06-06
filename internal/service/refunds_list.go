package service

import (
	"gitlab.ozon.dev/antonkraeww/homeworks/internal/domain/models"
	"gitlab.ozon.dev/antonkraeww/homeworks/internal/domain/requests"
)

// RefundsList returns list of refunds paginated
// optional pageN=<page number from the end>
// optional perPage=<number of orders per page>.
func (s *OrderService) RefundsList(req requests.RefundsListRequest) ([]models.Order, error) {
	orders, err := s.Repo.GetOrders(models.OrderFilter{
		Statuses:     []models.Status{models.Refunded},
		PageN:        req.PageN,
		PerPage:      req.PerPage,
		SortedByDate: true,
	})
	if err != nil {
		return []models.Order{}, err
	}

	return orders, nil
}
