package order

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.ozon.dev/antonkraeww/homeworks/hw-5/internal/models/domain/order"
	"gitlab.ozon.dev/antonkraeww/homeworks/hw-5/internal/models/domain/wrap"
	orderRepo "gitlab.ozon.dev/antonkraeww/homeworks/hw-5/internal/repository/order"
	"gitlab.ozon.dev/antonkraeww/homeworks/hw-5/tests/pg"
)

func TestGetOrderByID(t *testing.T) {
	var (
		ctx          = context.Background()
		ID    uint64 = 1
		badID uint64 = 2
		now          = time.Now().UTC()

		testWrap  = []wrap.Wrap{{Name: "box", Weight: 10, Cost: 10}}
		testOrder = []order.Order{
			{OrderID: ID, ClientID: 1, Weight: 1, Cost: 1, WrapType: "box", StoredUntil: now, StatusChanged: now},
		}
	)

	pg.DB.SetUp(t, "orders", "wrap")
	defer pg.DB.TearDown(t)
	pg.DB.FillWraps(testWrap)
	pg.DB.FillOrders(testOrder)
	repo := orderRepo.NewOrderRepository(pg.DB.ConnPool)

	t.Run("get_order", func(t *testing.T) {
		t.Parallel()

		resp, err := repo.GetOrderByID(ctx, ID)

		require.NoError(t, err)
		require.NotNil(t, resp)
		AssertEqualOrders(t, testOrder[0], *resp)
	})

	t.Run("not_found", func(t *testing.T) {
		t.Parallel()

		resp, err := repo.GetOrderByID(ctx, badID)

		require.NoError(t, err)
		assert.Nil(t, resp)
	})
}
