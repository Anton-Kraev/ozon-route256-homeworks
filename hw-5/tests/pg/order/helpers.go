package order

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"gitlab.ozon.dev/antonkraeww/homeworks/hw-5/internal/models/domain/order"
)

func AssertEqualOrders(t *testing.T, expected, actual order.Order) {
	assert.Equal(t, expected.OrderID, actual.OrderID)
	assert.Equal(t, expected.ClientID, actual.ClientID)
	assert.Equal(t, expected.Weight, actual.Weight)
	assert.Equal(t, expected.Cost, actual.Cost)
	assert.Equal(t, expected.Status, actual.Status)
	assert.Equal(t, expected.Hash, actual.Hash)
	assert.Equal(t, expected.WrapType, actual.WrapType)
	compareTimeWithTolerance(t, expected.StatusChanged, actual.StatusChanged)
	compareTimeWithTolerance(t, expected.StoredUntil, actual.StoredUntil)
}

func compareTimeWithTolerance(t *testing.T, expected time.Time, actual time.Time) {
	diff := actual.Sub(expected)

	assert.True(t, diff > -time.Millisecond && diff < time.Millisecond)
}
