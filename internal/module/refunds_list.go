package module

import (
	"gitlab.ozon.dev/antonkraeww/homeworks/internal/domain/models"
)

// RefundsList returns list of refunds paginated
// optional pageN=<page number from the end>
// optional perPage=<number of orders per page>.
func (m *OrderModule) RefundsList(pageN, perPage uint) ([]models.Order, error) {
	orders, err := m.Storage.GetOrders(models.OrderFilter{
		Statuses:     []models.Status{models.Refunded},
		PageN:        pageN,
		PerPage:      perPage,
		SortedByDate: true,
	})
	if err != nil {
		return []models.Order{}, err
	}

	return orders, nil
}
