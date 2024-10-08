//go:build integration

package pg

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"

	"github.com/Anton-Kraev/ozon-route256-homeworks/internal/middlewares"
	"github.com/Anton-Kraev/ozon-route256-homeworks/internal/models/domain/order"
	"github.com/Anton-Kraev/ozon-route256-homeworks/internal/models/domain/wrap"
	orderRepo "github.com/Anton-Kraev/ozon-route256-homeworks/internal/repository/order"
)

func TestAddOrder(t *testing.T) {
	var (
		now = time.Now().UTC()

		testWrap   = []wrap.Wrap{{Name: "box", MaxWeight: 10, Cost: 10}}
		testOrders = []order.Order{
			{
				OrderID:       1,
				ClientID:      1,
				Weight:        1,
				Cost:          1,
				Status:        order.Received,
				WrapType:      "box",
				StoredUntil:   now,
				StatusChanged: now,
				Hash:          "hash",
			},
			{
				OrderID:       2,
				ClientID:      2,
				Weight:        2,
				Cost:          2,
				Status:        order.Refunded,
				WrapType:      "badwrap",
				StoredUntil:   now,
				StatusChanged: now,
				Hash:          "hash",
			},
			{
				OrderID:       3,
				ClientID:      3,
				Weight:        3,
				Cost:          3,
				Status:        order.Returned,
				WrapType:      "",
				StoredUntil:   now,
				StatusChanged: now,
				Hash:          "hash",
			},
		}
	)

	DB.SetUp(t, "orders", "wrap")
	defer DB.TearDown(t)
	DB.fillWraps(testWrap)

	repo := orderRepo.NewOrderRepository(DB.ConnPool)
	txMw := middlewares.NewTransactionMiddleware(DB.ConnPool)

	t.Run("with_wrap", func(t *testing.T) {
		ctx := context.Background()
		res, err := txMw.CreateTransactionContext(
			ctx,
			pgx.TxOptions{AccessMode: pgx.ReadWrite, IsoLevel: pgx.RepeatableRead},
			[]string{},
			func(ctx context.Context, args []string) (string, error) {
				return "", repo.AddOrder(ctx, testOrders[0])
			},
		)

		require.NoError(t, err)
		require.Empty(t, res)

		records := DB.getAllOrders()
		require.NotEmpty(t, records)

		var actual order.Order

		for _, rec := range records {
			if rec.OrderID == 1 {
				actual = rec.ToDomain()
			}
		}

		AssertEqualOrders(t, testOrders[0], actual)
	})

	t.Run("bad_wrap", func(t *testing.T) {
		ctx := context.Background()
		_, err := txMw.CreateTransactionContext(
			ctx,
			pgx.TxOptions{AccessMode: pgx.ReadWrite, IsoLevel: pgx.RepeatableRead},
			[]string{},
			func(ctx context.Context, args []string) (string, error) {
				return "", repo.AddOrder(ctx, testOrders[1])
			},
		)

		require.Error(t, err)
	})

	t.Run("no_wrap", func(t *testing.T) {
		ctx := context.Background()
		res, err := txMw.CreateTransactionContext(
			ctx,
			pgx.TxOptions{AccessMode: pgx.ReadWrite, IsoLevel: pgx.RepeatableRead},
			[]string{},
			func(ctx context.Context, args []string) (string, error) {
				return "", repo.AddOrder(ctx, testOrders[2])
			},
		)

		require.NoError(t, err)
		require.Empty(t, res)

		records := DB.getAllOrders()
		require.NotEmpty(t, records)

		var actual order.Order

		for _, rec := range records {
			if rec.OrderID == 3 {
				actual = rec.ToDomain()
			}
		}

		AssertEqualOrders(t, testOrders[2], actual)
	})
}
