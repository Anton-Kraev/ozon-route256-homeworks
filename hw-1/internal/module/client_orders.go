package module

import (
	"sort"

	"gitlab.ozon.dev/antonkraeww/homeworks/hw-1/internal/domain/models"
)

// ClientOrders returns list of client orders
// optional lastN for get last orders, by default return all orders
// optional inStorage for get only orders from storage
func (m *OrderModule) ClientOrders(clientID uint64, lastN uint, onlyInStorage bool) ([]models.Order, error) {
	orders, err := m.Storage.ReadAll()
	if err != nil {
		return []models.Order{}, err
	}

	sort.Slice(orders, func(i, j int) bool {
		return orders[i].StatusChanged.After(orders[j].StatusChanged)
	})

	if lastN == 0 {
		lastN = uint(len(orders))
	}

	var clientOrders []models.Order

	for _, order := range orders {
		if lastN == 0 {
			break
		}

		inStorage := order.Status == models.Received || order.Status == models.Refunded
		if order.ClientID == clientID && (inStorage || !onlyInStorage) {
			clientOrders = append(clientOrders, order)
			lastN--
		}
	}

	return clientOrders, nil
}
