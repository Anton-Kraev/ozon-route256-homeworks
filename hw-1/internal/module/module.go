package module

import (
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"

	"gitlab.ozon.dev/antonkraeww/homeworks/hw-1/internal/models"
)

type Storage interface {
	AddOrder(newOrder models.Order) error
	ChangeOrders(changes map[uint64]models.Order) error
	FindOrder(orderID uint64) (*models.Order, error)
	ReadAll() ([]models.Order, error)
	RewriteAll(data []models.Order) error
}

type Deps struct {
	Storage Storage
}

type Module struct {
	Deps
	mu sync.Mutex
}

func NewModule(d Deps) Module {
	return Module{Deps: d}
}

// ReceiveOrder receives order from courier
func (m *Module) ReceiveOrder(orderID, clientID uint64, storedUntil time.Time) error {
	now := time.Now().UTC()
	if now.After(storedUntil) {
		return errors.New("retention time is in the past")
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
func (m *Module) ReturnOrder(orderID uint64) error {
	m.mu.Lock()
	order, err := m.Storage.FindOrder(orderID)
	m.mu.Unlock()
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	if order.StoredUntil.After(now) {
		return errors.New("retention period isn't expired yet")
	} else if order.Status == models.Returned {
		return errors.New("order has been already returned")
	} else if order.Status == models.Delivered {
		return errors.New("order has been delivered to client")
	}

	order.SetStatus(models.Returned, now)
	order.SetHash()

	m.mu.Lock()
	defer m.mu.Unlock()
	return m.Storage.ChangeOrders(map[uint64]models.Order{orderID: *order})
}

// DeliverOrders deliver list of orders to client
func (m *Module) DeliverOrders(ordersID []uint64) error {
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

	wg.Add(len(orders))

	for _, order := range orders {
		if _, ok := delivered[order.OrderID]; !ok {
			continue
		}

		now := time.Now().UTC()

		if prevDelivered.OrderID != 0 && order.ClientID != prevDelivered.ClientID {
			return fmt.Errorf(
				"orders with id %d and %d belong to different clients",
				order.ClientID, prevDelivered.ClientID,
			)
		}
		if order.Status != models.Received {
			return fmt.Errorf(
				"order with id %d has the %s status",
				order.OrderID, order.Status,
			)
		}
		if now.After(order.StoredUntil) {
			return fmt.Errorf(
				"retention period of order with id %d has expired",
				order.OrderID,
			)
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
			return fmt.Errorf("order with id %d not found", id)
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
func (m *Module) ClientOrders(clientID uint64, lastN uint, inStorage bool) ([]models.Order, error) {
	orders, err := m.Storage.ReadAll()
	if err != nil {
		return []models.Order{}, err
	}

	sort.Slice(orders, func(i, j int) bool {
		return orders[i].StatusChanged.After(orders[i].StatusChanged)
	})

	if lastN == 0 {
		lastN = uint(len(orders))
	}

	var clientOrders []models.Order
	for i := 0; i < len(orders) && lastN > 0; i++ {
		if orders[i].ClientID == clientID &&
			(!inStorage || orders[i].Status == models.Received || orders[i].Status == models.Refunded) {
			clientOrders = append(clientOrders, orders[i])
			lastN--
		}
	}

	return clientOrders, nil
}

// RefundOrder receives order refund from client
func (m *Module) RefundOrder(orderID, clientID uint64) error {
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
			return fmt.Errorf("order with id %d has been already refunded", order.OrderID)
		} else if order.Status != models.Delivered {
			return fmt.Errorf("order with id %d was not delivered to client yet", order.OrderID)
		}
		now := time.Now().UTC()
		if order.StatusChanged.Add(time.Hour * 48).Before(now) {
			return fmt.Errorf("more than two 2 days since order was delivered")
		}

		order.SetStatus(models.Refunded, now)
		order.SetHash()

		return m.Storage.ChangeOrders(map[uint64]models.Order{orderID: order})
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	return fmt.Errorf("order of client %d with id %d not found", clientID, orderID)
}

// RefundsList returns list of refunds paginated
// optional pageN=<page number from the end>
// optional perPage=<number of orders per page>
func (m *Module) RefundsList(pageN, perPage uint) ([]models.Order, error) {
	orders, err := m.Storage.ReadAll()
	if err != nil {
		return []models.Order{}, err
	}

	sort.Slice(orders, func(i, j int) bool {
		return orders[i].StatusChanged.After(orders[i].StatusChanged)
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
	return refunds[pageN*perPage : max(uint(len(refunds)), (pageN+1)*perPage)], nil
}
