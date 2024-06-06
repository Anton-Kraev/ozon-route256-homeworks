package repository

import (
	"math"
	"sort"

	"gitlab.ozon.dev/antonkraeww/homeworks/internal/domain/models"
)

// GetOrders returns orders that matches specified filter.
func (r OrderRepository) GetOrders(filter models.OrderFilter) ([]models.Order, error) {
	orders, err := r.readAll()
	if err != nil {
		return []models.Order{}, err
	}

	filter.Init()

	// preventing incorrect work if overflow of uint in filter.PageN*filter.PerPage
	overflow := uint(math.Ceil(float64(math.MaxUint)/float64(filter.PerPage))) <= filter.PageN
	if filter.PageN*filter.PerPage >= uint(len(orders)) || overflow {
		return []models.Order{}, nil
	}

	if filter.SortedByDate {
		sort.Slice(orders, func(i, j int) bool {
			return orders[i].StatusChanged.After(orders[j].StatusChanged)
		})
	}

	var filteredOrders []models.Order
	for _, order := range orders {
		if order.MatchesFilter(filter) {
			filteredOrders = append(filteredOrders, order)
		}

		if uint(len(filteredOrders)) >= filter.PerPage {
			break
		}
	}

	return filteredOrders, nil
}
