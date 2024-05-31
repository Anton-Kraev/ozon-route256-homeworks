package storage

import (
	"encoding/json"
	"errors"
	"os"

	"gitlab.ozon.dev/antonkraeww/homeworks/hw-1/internal/domain/models"
)

type OrderStorage struct {
	fileName string
}

func NewOrderStorage(fileName string) OrderStorage {
	return OrderStorage{fileName: fileName}
}

// AddOrder adds new order to end of storage (if passed ID param is unique)
func (s OrderStorage) AddOrder(newOrder models.Order) error {
	orders, err := s.ReadAll()
	if err != nil {
		return err
	}

	for _, order := range orders {
		if order.OrderID == newOrder.OrderID {
			return errors.New("order ID must be unique")
		}
	}

	return s.RewriteAll(append(orders, newOrder))
}

// ChangeOrders changes orders data in storage, key=<order id to change> value=<new order data>
func (s OrderStorage) ChangeOrders(changes map[uint64]models.Order) error {
	orders, err := s.ReadAll()
	if err != nil {
		return err
	}

	for i, order := range orders {
		if _, ok := changes[order.OrderID]; ok {
			orders[i] = changes[order.OrderID]
		}
	}

	return s.RewriteAll(orders)
}

// FindOrder find order with specified orderID in storage
func (s OrderStorage) FindOrder(orderID uint64) (*models.Order, error) {
	orders, err := s.ReadAll()
	if err != nil {
		return nil, err
	}

	for _, order := range orders {
		if order.OrderID == orderID {
			return &order, nil
		}
	}

	return nil, errors.New("order not found")
}

// ReadAll return all orders
func (s OrderStorage) ReadAll() ([]models.Order, error) {
	if _, err := os.Stat(s.fileName); errors.Is(err, os.ErrNotExist) {
		f, errCreate := os.Create(s.fileName)
		if errCreate != nil {
			return []models.Order{}, errCreate
		}

		if errClose := f.Close(); errClose != nil {
			return []models.Order{}, errClose
		}
		return []models.Order{}, nil
	}

	bytes, errRead := os.ReadFile(s.fileName)
	if errRead != nil {
		return []models.Order{}, errRead
	}

	var records []orderRecord
	if errUnmarshal := json.Unmarshal(bytes, &records); errUnmarshal != nil {
		return []models.Order{}, errUnmarshal
	}

	var data []models.Order
	for _, record := range records {
		data = append(data, record.toDomain())
	}

	return data, nil
}

// RewriteAll rewrites storage with specified data
func (s OrderStorage) RewriteAll(data []models.Order) error {
	var orders []orderRecord
	for _, order := range data {
		orders = append(orders, toRecord(order))
	}

	bytes, err := json.MarshalIndent(orders, "  ", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.fileName, bytes, 0644)
}
