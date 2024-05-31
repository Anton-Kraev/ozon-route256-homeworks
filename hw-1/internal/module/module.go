package module

import (
	"sort"
	"sync"
	"time"

	domainErrors "gitlab.ozon.dev/antonkraeww/homeworks/hw-1/internal/domain/errors"
	"gitlab.ozon.dev/antonkraeww/homeworks/hw-1/internal/domain/models"
)

type orderStorage interface {
	AddOrder(newOrder models.Order) error
	ChangeOrders(changes map[uint64]models.Order) error
	FindOrder(orderID uint64) (*models.Order, error)
	ReadAll() ([]models.Order, error)
	RewriteAll(data []models.Order) error
}

type OrderModule struct {
	Storage orderStorage
	mu      sync.Mutex
}

func NewOrderModule(storage orderStorage) OrderModule {
	return OrderModule{Storage: storage}
}

// ReceiveOrder receives order from courier
func (m *OrderModule) ReceiveOrder(orderID, clientID uint64, storedUntil time.Time) error {
	now := time.Now().UTC()
	if now.After(storedUntil) {
		return domainErrors.ErrRetentionTimeInPast
	}

	newOrder := models.Order{
		OrderID:       orderID,
		ClientID:      clientID,
		StoredUntil:   storedUntil,
		Status:        models.Received,
		StatusChanged: now,
	}
	newOrder.SetHash()

	m.mu.Lock()
	defer m.mu.Unlock()
	return m.Storage.AddOrder(newOrder)
}

// ReturnOrder returns order to courier
func (m *OrderModule) ReturnOrder(orderID uint64) error {
	m.mu.Lock()
	order, err := m.Storage.FindOrder(orderID)
	m.mu.Unlock()
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	if order.StoredUntil.After(now) {
		return domainErrors.ErrRetentionPeriodNotExpiredYet
	}
	if order.Status == models.Returned {
		return domainErrors.ErrOrderAlreadyReturned
	}
	if order.Status == models.Delivered {
		return domainErrors.ErrOrderDelivered
	}

	order.SetStatus(models.Returned, now)
	order.SetHash()

	m.mu.Lock()
	defer m.mu.Unlock()
	return m.Storage.ChangeOrders(map[uint64]models.Order{orderID: *order})
}

// DeliverOrders deliver list of orders to client
func (m *OrderModule) DeliverOrders(ordersID []uint64) error {
	delivered := make(map[uint64]*models.Order)
	for _, orderID := range ordersID {
		delivered[orderID] = nil
	}

	m.mu.Lock()
	orders, err := m.Storage.ReadAll()
	m.mu.Unlock()
	if err != nil {
		return err
	}

	var (
		prevDelivered models.Order
		wg            sync.WaitGroup
	)

	wg.Add(len(ordersID))

	for _, order := range orders {
		if _, ok := delivered[order.OrderID]; !ok {
			continue
		}

		now := time.Now().UTC()

		if prevDelivered.OrderID != 0 && order.ClientID != prevDelivered.ClientID {
			return domainErrors.ErrDifferentClientOrders(order.ClientID, prevDelivered.ClientID)
		}
		if order.Status != models.Received {
			return domainErrors.ErrUnexpectedOrderStatus(order.OrderID, order.Status)
		}
		if now.After(order.StoredUntil) {
			return domainErrors.ErrRetentionPeriodExpired(order.OrderID)
		}

		prevDelivered = order

		go func(order models.Order) {
			defer wg.Done()
			order.SetStatus(models.Delivered, now)
			order.SetHash()
			delivered[order.OrderID] = &order
		}(order)
	}

	wg.Wait()

	changes := make(map[uint64]models.Order)
	for id, order := range delivered {
		if order == nil {
			return domainErrors.ErrOrderNotFound(id)
		}
		changes[id] = *order
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	return m.Storage.ChangeOrders(changes)
}

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

// RefundOrder receives order refund from client
func (m *OrderModule) RefundOrder(orderID, clientID uint64) error {
	m.mu.Lock()
	orders, err := m.Storage.ReadAll()
	m.mu.Unlock()
	if err != nil {
		return err
	}

	for _, order := range orders {
		if order.OrderID != orderID || order.ClientID != clientID {
			continue
		}

		if order.Status == models.Refunded {
			return domainErrors.ErrOrderAlreadyRefunded
		}
		if order.Status != models.Delivered {
			return domainErrors.ErrOrderNotDeliveredYet
		}
		now := time.Now().UTC()
		if order.StatusChanged.Add(time.Hour * 48).Before(now) {
			return domainErrors.ErrOrderDeliveredLongAgo
		}

		order.SetStatus(models.Refunded, now)
		order.SetHash()

		return m.Storage.ChangeOrders(map[uint64]models.Order{orderID: order})
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	return domainErrors.ErrOrderNotFound(orderID)
}

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
