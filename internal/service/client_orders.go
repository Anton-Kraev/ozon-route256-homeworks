package service

import (
	"gitlab.ozon.dev/antonkraeww/homeworks/internal/domain/models"
)

// ClientOrders returns list of client orders
// optional lastN for get last orders, by default return all orders
// optional inStorage for get only orders from storage.
func (s *OrderService) ClientOrders(clientID uint64, lastN uint, onlyInStorage bool) ([]models.Order, error) {
	filter := models.OrderFilter{
		ClientsID:    []uint64{clientID},
		PerPage:      lastN,
		SortedByDate: true,
	}
	if onlyInStorage {
		filter.Statuses = []models.Status{models.Received, models.Refunded}
	}

	return s.Repo.GetOrders(filter)
}
