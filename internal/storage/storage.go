package storage

import (
	"encoding/json"
	"errors"
	"os"

	"gitlab.ozon.dev/antonkraeww/homeworks/internal/domain/models"
)

type OrderStorage struct {
	fileName string
}

func NewOrderStorage(fileName string) OrderStorage {
	return OrderStorage{fileName: fileName}
}

// readAll return all orders.
func (s OrderStorage) readAll() ([]models.Order, error) {
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

// rewriteAll rewrites storage with specified data.
func (s OrderStorage) rewriteAll(data []models.Order) error {
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
