package module

import (
	"errors"
	"fmt"
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

func (m *Module) ClientOrders(clientID uint64, lastN uint, inStorage bool) ([]models.Order, error) {
	orders, err := m.Storage.ReadAll()
	if err != nil {
		return []models.Order{}, err
	}

	if lastN == 0 {
		lastN = uint(len(orders))
	}

	var clientOrders []models.Order
	for i := len(orders) - 1; i >= 0 && lastN > 0; i-- {
		if orders[i].ClientID == clientID &&
			(!inStorage || orders[i].Status == models.Received || orders[i].Status == models.Refunded) {
			clientOrders = append(clientOrders, orders[i])
		}
	}

	return clientOrders, nil
}

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
			return fmt.Errorf("more than two 2 days since order was deliverec")
		}

		order.SetStatus(models.Refunded, now)
		order.SetHash()

		return m.Storage.ChangeOrders(map[uint64]models.Order{orderID: order})
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	return fmt.Errorf("order of client %d with id %d not found", clientID, orderID)
}

func (m *Module) RefundsList(pageN, perPage uint) ([]models.Order, error) {
	orders, err := m.Storage.ReadAll()
	if err != nil {
		return []models.Order{}, err
	}

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
	return refunds[pageN*perPage : (pageN+1)*perPage], err
}
