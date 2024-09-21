package order

import (
	"context"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	errsdomain "gitlab.ozon.dev/antonkraeww/ozon-route256-homeworks/internal/models/domain/errors"
	"gitlab.ozon.dev/antonkraeww/ozon-route256-homeworks/internal/models/domain/order"
	service "gitlab.ozon.dev/antonkraeww/ozon-route256-homeworks/internal/service/order"
)

func TestOrderController_DeliverOrders(t *testing.T) {
	t.Parallel()

	type fields struct {
		orderRepo *service.MockorderRepository
		wrapRepo  *service.MockwrapRepository
		hashes    *service.MockhashGenerator
	}

	var (
		ctx = context.Background()
	)

	tests := []struct {
		name    string
		input   string
		fields  fields
		mockFn  func(f fields)
		wantErr error
	}{
		{
			name:    "no_orders",
			input:   "deliver",
			mockFn:  func(f fields) {},
			wantErr: errParseArgs,
		},
		{
			name:    "bad_format",
			input:   "deliver --orders=1;2;3",
			mockFn:  func(f fields) {},
			wantErr: errParseArgs,
		},
		{
			name:    "negative_orderid",
			input:   "deliver --orders=-1",
			mockFn:  func(f fields) {},
			wantErr: errParseArgs,
		},
		{
			name:  "successful_parse",
			input: "deliver --orders=1,2",
			mockFn: func(f fields) {
				f.orderRepo.EXPECT().GetOrdersByIDs(gomock.Any(), []uint64{1, 2}).Return([]order.Order{}, nil)
			},
			wantErr: errsdomain.ErrOrderNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			tt.fields = fields{
				orderRepo: service.NewMockorderRepository(ctrl),
				wrapRepo:  service.NewMockwrapRepository(ctrl),
				hashes:    service.NewMockhashGenerator(ctrl),
			}
			tt.mockFn(tt.fields)

			srvc := service.NewOrderService(tt.fields.orderRepo, tt.fields.wrapRepo, tt.fields.hashes)
			controller := NewOrderController(srvc)

			res, err := controller.DeliverOrders(ctx, strings.Split(tt.input, " ")[1:])

			require.ErrorIs(t, err, tt.wantErr)
			assert.Empty(t, res)
		})
	}
}
