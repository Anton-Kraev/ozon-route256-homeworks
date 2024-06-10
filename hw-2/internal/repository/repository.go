package repository

import (
	"encoding/json"
	"errors"
	"os"

	models "gitlab.ozon.dev/antonkraeww/homeworks/hw-2/internal/models/domain/order"
)

type OrderRepository struct {
	fileName string
}

func NewOrderRepository(fileName string) OrderRepository {
	return OrderRepository{fileName: fileName}
}

// readAll return all orders.
func (r OrderRepository) readAll() ([]models.Order, error) {
	if _, err := os.Stat(r.fileName); errors.Is(err, os.ErrNotExist) {
		f, errCreate := os.Create(r.fileName)

		if errCreate != nil {
			return []models.Order{}, errCreate
		}
		if errClose := f.Close(); errClose != nil {
			return []models.Order{}, errClose
		}

		return []models.Order{}, nil
	}

	bytes, errRead := os.ReadFile(r.fileName)
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

// rewriteAll rewrites storage with specified data.
func (r OrderRepository) rewriteAll(data []models.Order) error {
	var orders []orderRecord
	for _, order := range data {
		orders = append(orders, toRecord(order))
	}

	bytes, err := json.MarshalIndent(orders, "  ", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(r.fileName, bytes, 0644)
}
