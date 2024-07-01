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

func TestGetOrdersByFilter(t *testing.T) {
	var (
		ctx = context.Background()
		now = time.Now().UTC()

		testWrap   = []wrap.Wrap{{Name: "box", Weight: 10, Cost: 10}}
		testOrders = []order.Order{
			{OrderID: 1, ClientID: 1, Weight: 1, Cost: 1, WrapType: "box", Status: order.Received, StatusChanged: now.Add(time.Hour * 4), StoredUntil: now},
			{OrderID: 2, ClientID: 1, Weight: 2, Cost: 2, WrapType: "box", Status: order.Refunded, StatusChanged: now.Add(time.Hour * 3), StoredUntil: now},
			{OrderID: 3, ClientID: 1, Weight: 3, Cost: 3, WrapType: "box", Status: order.Delivered, StatusChanged: now.Add(time.Hour * 2), StoredUntil: now},
			{OrderID: 4, ClientID: 2, Weight: 4, Cost: 4, WrapType: "box", Status: order.Returned, StatusChanged: now.Add(time.Hour * 1), StoredUntil: now},
		}
	)

	pg.DB.SetUp(t, "orders", "wrap")
	defer pg.DB.TearDown(t)
	pg.DB.FillWraps(testWrap)
	pg.DB.FillOrders(testOrders)
	repo := orderRepo.NewOrderRepository(pg.DB.ConnPool)

	t.Run("complex_filter", func(t *testing.T) {
		t.Parallel()

		resp, err := repo.GetOrdersByFilter(ctx, order.Filter{
			ClientID:     1,
			Statuses:     []order.Status{order.Refunded, order.Delivered},
			PageN:        0,
			PerPage:      2,
			SortedByDate: true,
		})

		require.NoError(t, err)
		require.Len(t, resp, 2)
		AssertEqualOrders(t, testOrders[1], resp[0])
		AssertEqualOrders(t, testOrders[2], resp[1])
	})

	t.Run("unknown_client", func(t *testing.T) {
		t.Parallel()

		resp, err := repo.GetOrdersByFilter(ctx, order.Filter{ClientID: 3})

		require.NoError(t, err)
		require.Empty(t, resp)
	})

	t.Run("last_record", func(t *testing.T) {
		t.Parallel()

		resp, err := repo.GetOrdersByFilter(ctx, order.Filter{PageN: 1, PerPage: 3, SortedByDate: true})

		require.NoError(t, err)
		require.Len(t, resp, 1)
		AssertEqualOrders(t, testOrders[3], resp[0])
	})

	t.Run("get_all", func(t *testing.T) {
		t.Parallel()

		resp, err := repo.GetOrdersByFilter(ctx, order.Filter{SortedByDate: true})

		require.NoError(t, err)
		require.Len(t, resp, 4)
		AssertEqualOrders(t, testOrders[0], resp[0])
		AssertEqualOrders(t, testOrders[1], resp[1])
		AssertEqualOrders(t, testOrders[2], resp[2])
		AssertEqualOrders(t, testOrders[3], resp[3])
	})

	t.Run("bad_filter", func(t *testing.T) {
		t.Parallel()

		resp, err := repo.GetOrdersByFilter(ctx, order.Filter{PageN: 1})

		require.NoError(t, err)
		require.Empty(t, resp)
	})
}
