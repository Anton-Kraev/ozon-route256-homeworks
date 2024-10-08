package errors

import (
	"errors"
	"fmt"
)

var (
	ErrRetentionTimeInPast          = errors.New("retention time is in the past")
	ErrOrderAlreadyReturned         = errors.New("order has been already returned")
	ErrOrderAlreadyRefunded         = errors.New("order has been already refunded")
	ErrOrderDelivered               = errors.New("order has been delivered to client")
	ErrOrderNotDeliveredYet         = errors.New("order was not delivered to client yet")
	ErrRetentionPeriodNotExpiredYet = errors.New("retention period isn't expired yet")
	ErrOrderDeliveredLongAgo        = errors.New("more than two 2 days since order was delivered")
	ErrDifferentClients             = errors.New("orders belong to different clients")
	ErrUnexpectedOrderStatus        = errors.New("unexpected order status")
	ErrRetentionPeriodExpired       = errors.New("retention period is expired")
	ErrOrderNotFound                = errors.New("order not found")
	ErrOrderIDNotUnique             = errors.New("order with same ID already exist")
	ErrUnknownOrderWrapType         = errors.New("unknown order wrap type")
	ErrOrderWeightExceedsLimit      = errors.New("order weight exceeds permissible weight for this type of wrap")

	ErrWrapNotFound      = errors.New("wrap not found")
	ErrWrapAlreadyExists = errors.New("wrap already exists")
)

func ErrorDifferentClients(clientID, anotherClientID uint64) error {
	return fmt.Errorf("orders with id %d and %d belong to different clients", clientID, anotherClientID)
}

func ErrorUnexpectedOrderStatus(orderID uint64, status string) error {
	return fmt.Errorf("order with id %d has the %s status", orderID, status)
}

func ErrorRetentionPeriodExpired(orderID uint64) error {
	return fmt.Errorf("retention period of order with id %d has expired", orderID)
}

func ErrorOrderNotFound(orderID uint64) error {
	return fmt.Errorf("order with id %d not found", orderID)
}
