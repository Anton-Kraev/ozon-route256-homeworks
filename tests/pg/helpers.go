package pg

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/Anton-Kraev/ozon-route256-homeworks/internal/models/domain/order"
	"github.com/Anton-Kraev/ozon-route256-homeworks/internal/models/domain/wrap"
)

func AssertEqualWraps(t *testing.T, expected, actual wrap.Wrap) {
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.MaxWeight, actual.MaxWeight)
	assert.Equal(t, expected.Cost, actual.Cost)
}

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
