package order

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"gitlab.ozon.dev/antonkraeww/homeworks/hw-5/internal/models/domain/order"
	"gitlab.ozon.dev/antonkraeww/homeworks/hw-5/internal/models/domain/wrap"
	orderRepo "gitlab.ozon.dev/antonkraeww/homeworks/hw-5/internal/repository/order"
	"gitlab.ozon.dev/antonkraeww/homeworks/hw-5/tests/pg"
)

func TestGetOrdersByIDs(t *testing.T) {
	var (
		ctx          = context.Background()
		IDs          = []uint64{1, 2}
		badID uint64 = 3
		now          = time.Now().UTC()

		testWrap   = []wrap.Wrap{{Name: "box", Weight: 10, Cost: 10}}
		testOrders = []order.Order{
			{OrderID: IDs[0], ClientID: 1, Weight: 1, Cost: 1, WrapType: "box", StoredUntil: now, StatusChanged: now},
			{OrderID: IDs[1], ClientID: 2, Weight: 2, Cost: 2, WrapType: "box", StoredUntil: now, StatusChanged: now},
		}
	)

	pg.DB.SetUp(t, "orders", "wrap")
	defer pg.DB.TearDown(t)
	pg.DB.FillWraps(testWrap)
	pg.DB.FillOrders(testOrders)
	repo := orderRepo.NewOrderRepository(pg.DB.ConnPool)

	t.Run("get_orders", func(t *testing.T) {
		t.Parallel()

		resp, err := repo.GetOrdersByIDs(ctx, IDs)

		require.NoError(t, err)
		require.Len(t, resp, len(testOrders))
		AssertEqualOrders(t, testOrders[0], resp[0])
		AssertEqualOrders(t, testOrders[1], resp[1])
	})

	t.Run("one_not_found", func(t *testing.T) {
		t.Parallel()

		resp, err := repo.GetOrdersByIDs(ctx, append(IDs, badID))

		require.NoError(t, err)
		require.Len(t, resp, len(testOrders))
		AssertEqualOrders(t, testOrders[0], resp[0])
		AssertEqualOrders(t, testOrders[1], resp[1])
	})
}
