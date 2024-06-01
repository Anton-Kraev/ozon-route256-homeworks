package errsdomain

import (
	"errors"
	"fmt"

	"gitlab.ozon.dev/antonkraeww/homeworks/hw-1/internal/domain/models"
)

var (
	ErrRetentionTimeInPast          = errors.New("retention time is in the past")
	ErrOrderAlreadyReturned         = errors.New("order has been already returned")
	ErrOrderAlreadyRefunded         = errors.New("order has been already refunded")
	ErrOrderDelivered               = errors.New("order has been delivered to client")
	ErrOrderNotDeliveredYet         = errors.New("order was not delivered to client yet")
	ErrRetentionPeriodNotExpiredYet = errors.New("retention period isn't expired yet")
	ErrOrderDeliveredLongAgo        = errors.New("more than two 2 days since order was delivered")
)

func ErrDifferentClientOrders(clientID, anotherClientID uint64) error {
	return fmt.Errorf("orders with id %d and %d belong to different clients", clientID, anotherClientID)
}

func ErrUnexpectedOrderStatus(orderID uint64, status models.Status) error {
	return fmt.Errorf("order with id %d has the %s status", orderID, status)
}

func ErrRetentionPeriodExpired(orderID uint64) error {
	return fmt.Errorf("retention period of order with id %d has expired", orderID)
}

func ErrOrderNotFound(orderID uint64) error {
	return fmt.Errorf("order with id %d not found", orderID)
}

func ErrOrderIDNotUnique(orderID uint64) error {
	return fmt.Errorf("order with ID %d already exist", orderID)
}
