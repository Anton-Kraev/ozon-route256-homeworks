package storage

import (
	"encoding/json"
	"errors"
	"os"

	"gitlab.ozon.dev/antonkraeww/homeworks/hw-1/internal/models"
)

type Storage struct {
	fileName string
}

func NewStorage(fileName string) Storage {
	return Storage{fileName: fileName}
}

func (s Storage) AddOrder(newOrder models.Order) error { // TODO: add param sort by date
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

func (s Storage) ChangeOrders(changes map[int64]models.Order) error {
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

func (s Storage) RemoveOrder(orderID int64) error {
	orders, err := s.ReadAll()
	if err != nil {
		return err
	}
	for i, order := range orders {
		if order.OrderID == orderID {
			orders[i] = orders[len(orders)-1]
			return s.RewriteAll(orders)
		}
	}
	return errors.New("order not found")
}

func (s Storage) FindOrder(orderID int64) (*models.Order, error) {
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

func (s Storage) ReadAll() ([]models.Order, error) {
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
	b, errRead := os.ReadFile(s.fileName)
	if errRead != nil {
		return []models.Order{}, errRead
	}

	var records []orderRecord
	if errUnmarshal := json.Unmarshal(b, &records); errUnmarshal != nil {
		return []models.Order{}, errUnmarshal
	}
	var data []models.Order
	for _, record := range records {
		data = append(data, record.toDomain())
	}
	return data, nil
}

func (s Storage) RewriteAll(data []models.Order) error {
	var orders []orderRecord
	for _, order := range data {
		orders = append(orders, toRecord(order))
	}
	b, err := json.MarshalIndent(orders, "  ", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.fileName, b, 0644)
}
