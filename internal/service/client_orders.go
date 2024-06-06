package service

import (
	"gitlab.ozon.dev/antonkraeww/homeworks/internal/domain/models"
	"gitlab.ozon.dev/antonkraeww/homeworks/internal/domain/requests"
)

// ClientOrders returns list of client orders
// optional lastN for get last orders, by default return all orders
// optional inStorage for get only orders from storage.
func (s *OrderService) ClientOrders(req requests.ClientOrdersRequest) ([]models.Order, error) {
	filter := models.OrderFilter{
		ClientsID:    []uint64{req.ClientID},
		PerPage:      req.LastN,
		SortedByDate: true,
	}
	if req.InStorage {
		filter.Statuses = []models.Status{models.Received, models.Refunded}
	}

	return s.Repo.GetOrders(filter)
}
