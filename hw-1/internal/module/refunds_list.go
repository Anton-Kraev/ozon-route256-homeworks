package module

import (
	"sort"

	"gitlab.ozon.dev/antonkraeww/homeworks/hw-1/internal/domain/models"
)

// RefundsList returns list of refunds paginated
// optional pageN=<page number from the end>
// optional perPage=<number of orders per page>
func (m *OrderModule) RefundsList(pageN, perPage uint) ([]models.Order, error) {
	orders, err := m.Storage.ReadAll()
	if err != nil {
		return []models.Order{}, err
	}

	sort.Slice(orders, func(i, j int) bool {
		return orders[i].StatusChanged.After(orders[j].StatusChanged)
	})

	var refunds []models.Order
	for _, order := range orders {
		if order.Status == models.Refunded {
			refunds = append(refunds, order)
		}
	}

	if perPage > uint(len(refunds)) || perPage == 0 {
		perPage = uint(len(refunds))
	}

	if pageN*perPage >= uint(len(refunds)) {
		return []models.Order{}, nil
	}
	return refunds[pageN*perPage : min(uint(len(refunds)), (pageN+1)*perPage)], nil
}
