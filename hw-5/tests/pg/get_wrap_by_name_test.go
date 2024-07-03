//go:build integration

package pg

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.ozon.dev/antonkraeww/homeworks/hw-5/internal/models/domain/wrap"
	wrapRepo "gitlab.ozon.dev/antonkraeww/homeworks/hw-5/internal/repository/wrap"
)

func TestGetWrapByName(t *testing.T) {
	var (
		ctx     = context.Background()
		names   = []string{"name1", "name2"}
		badName = "unknown"

		testWraps = []wrap.Wrap{
			{
				Name:      names[0],
				MaxWeight: 1000,
				Cost:      1000,
			},
			{
				Name:      names[1],
				MaxWeight: 2000,
				Cost:      2000,
			},
		}
	)

	DB.SetUp(t, "orders", "wrap")
	defer DB.TearDown(t)
	DB.fillWraps(testWraps)

	repo := wrapRepo.NewWrapRepository(DB.ConnPool)

	t.Run("get_wrap", func(t *testing.T) {
		resp, err := repo.GetWrapByName(ctx, names[0])

		require.NoError(t, err)
		require.NotNil(t, resp)
		AssertEqualWraps(t, testWraps[0], *resp)
	})

	t.Run("not_found", func(t *testing.T) {
		resp, err := repo.GetWrapByName(ctx, badName)

		require.NoError(t, err)
		assert.Nil(t, resp)
	})
}
