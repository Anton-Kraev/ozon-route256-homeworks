package wrap

import (
	"context"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	errsdomain "gitlab.ozon.dev/antonkraeww/homeworks/hw-5/internal/models/domain/errors"
	"gitlab.ozon.dev/antonkraeww/homeworks/hw-5/internal/models/domain/wrap"
	service "gitlab.ozon.dev/antonkraeww/homeworks/hw-5/internal/service/wrap"
)

func TestWrapController_AddWrap(t *testing.T) {
	t.Parallel()

	var (
		ctx     = context.Background()
		oldWrap = "box"
	)

	tests := []struct {
		name    string
		input   string
		mockFn  func(repo *service.MockwrapRepository)
		wantErr error
	}{
		{
			"negative_weight",
			"addwrap --wrap=box --weight=-1 --cost=1",
			func(repo *service.MockwrapRepository) {},
			errParseArgs,
		},
		{
			"negative_cost",
			"addwrap --name=box --weight=1 --cost=-1",
			func(repo *service.MockwrapRepository) {},
			errParseArgs,
		},
		{
			"empty_name",
			"addwrap --weight=1 --cost=1",
			func(repo *service.MockwrapRepository) {},
			errParseArgs,
		},
		{
			"add_existing",
			"addwrap --name=box --weight=100 --cost=100",
			func(repo *service.MockwrapRepository) {
				repo.EXPECT().GetWrapByName(gomock.Any(), oldWrap).Return(&wrap.Wrap{Name: oldWrap}, nil)
			},
			errsdomain.ErrWrapAlreadyExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockRepo := service.NewMockwrapRepository(gomock.NewController(t))
			tt.mockFn(mockRepo)
			srvc := service.NewWrapService(mockRepo)
			ctrl := NewWrapController(srvc)

			res, err := ctrl.AddWrap(ctx, strings.Split(tt.input, " ")[1:])

			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, "", res)
		})
	}
}
