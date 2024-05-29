package module

import (
	"errors"
	"fmt"
	"time"

	"gitlab.ozon.dev/antonkraeww/homeworks/hw-1/internal/models"
)

type Storage interface {
	AddOrder(newOrder models.Order) error
	ChangeOrders(changes map[uint64]models.Order) error
	RemoveOrder(orderID uint64) error
	FindOrder(orderID uint64) (*models.Order, error)
	ReadAll() ([]models.Order, error)
	RewriteAll(data []models.Order) error
}

type Deps struct {
	Storage Storage
}

type Module struct {
	Deps
}

func NewModule(d Deps) Module {
	return Module{Deps: d}
}

func (m Module) ReceiveOrder(orderID, clientID uint64, storedUntil time.Time) error {
	now := time.Now().UTC()
	if now.After(storedUntil) {
		return errors.New("retention time is in the past")
	}

	return m.Storage.AddOrder(models.Order{
		OrderID:       orderID,
		ClientID:      clientID,
		StoredUntil:   storedUntil,
		Status:        models.Received,
		StatusChanged: now,
	})
}

func (m Module) ReturnOrder(orderId uint64) error {
	order, err := m.Storage.FindOrder(orderId)
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

	order.Status = models.Returned
	order.StatusChanged = now
	return m.Storage.ChangeOrders(map[uint64]models.Order{orderId: *order})
}

func (m Module) DeliverOrders(ordersID []uint64) error {
	delivered := make(map[uint64]*models.Order)
	for _, orderID := range ordersID {
		delivered[orderID] = nil
	}

	orders, err := m.Storage.ReadAll()
	if err != nil {
		return err
	}

	var prevDelivered *models.Order
	for _, order := range orders {
		now := time.Now().UTC()
		if prevDelivered != nil {
			var errDifferentClient, errNotReceived, errExpired error
			if order.ClientID != prevDelivered.ClientID {
				errDifferentClient = fmt.Errorf(
					"orders with id %d and %d belong to different clients",
					order.ClientID, prevDelivered.ClientID,
				)
			} else if order.Status != models.Received {
				errNotReceived = fmt.Errorf(
					"order with id %d has the %s status",
					order.OrderID, order.Status,
				)
			} else if now.After(order.StoredUntil) {
				errExpired = fmt.Errorf(
					"retention period of order with id %d has expired",
					order.OrderID,
				)
			}
			if errDeliver := errors.Join(errExpired, errNotReceived, errDifferentClient); errDeliver != nil {
				return errDeliver
			}
		}

		order.Status = models.Delivered
		order.StatusChanged = time.Now().UTC()
		delivered[order.OrderID] = &order
		prevDelivered = &order
	}

	changes := make(map[uint64]models.Order)
	for id, order := range delivered {
		if order == nil {
			return fmt.Errorf("order with id %d not found", id)
		}
		changes[id] = *order
	}

	return m.Storage.ChangeOrders(changes)
}

func (m Module) ClientOrders(clientID uint64, lastN uint, inStorage bool) ([]models.Order, error) {
	orders, err := m.Storage.ReadAll()
	if err != nil {
		return []models.Order{}, err
	}

	if lastN == 0 {
		lastN = uint(len(orders))
	}
	var clientOrders []models.Order
	for i := len(orders) - 1; i >= 0 && lastN > 0; i-- {
		if orders[i].ClientID == clientID && (!inStorage || orders[i].Status == models.Received || orders[i].Status == models.Refunded) {
			clientOrders = append(clientOrders, orders[i])
		}
	}

	return clientOrders, nil
}

func (m Module) RefundOrder(orderID, clientID uint64) error {
	orders, err := m.Storage.ReadAll()
	if err != nil {
		return err
	}

	for _, order := range orders {
		if order.OrderID == orderID && order.ClientID == clientID {
			if order.Status == models.Refunded {
				return fmt.Errorf("order with id %d has been already refunded", order.OrderID)
			} else if order.Status != models.Delivered {
				return fmt.Errorf("order with id %d was not delivered to client yet", order.OrderID)
			}
			now := time.Now().UTC()
			if now.Add(time.Hour * 48).After(order.StatusChanged) {
				return fmt.Errorf("more than two 2 days since order was deliverec")
			}

			order.Status = models.Refunded
			order.StatusChanged = now
			return m.Storage.ChangeOrders(map[uint64]models.Order{orderID: order})
		}
	}

	return fmt.Errorf("order of client %d with id %d not found", clientID, orderID)
}

func (m Module) RefundsList(pageN, perPage uint) ([]models.Order, error) {
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

	if pageN*perPage >= uint(len(refunds)) {
		return []models.Order{}, nil
	}
	return refunds[pageN*perPage : min(uint(len(refunds)), (pageN+1)*perPage)], err
}
