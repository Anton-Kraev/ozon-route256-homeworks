package order

import (
	"context"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.ozon.dev/antonkraeww/homeworks/hw-6/internal/models/domain/order"
	service "gitlab.ozon.dev/antonkraeww/homeworks/hw-6/internal/service/order"
)

func TestOrderController_ClientOrders(t *testing.T) {
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
		wantRes string
		wantErr error
	}{
		{
			name:    "no_clientid",
			input:   "olist",
			mockFn:  func(f fields) {},
			wantRes: "",
			wantErr: errParseArgs,
		},
		{
			name:    "bad_instorage",
			input:   "olist --clientID=12345 --inStorage=100",
			mockFn:  func(f fields) {},
			wantRes: "",
			wantErr: errParseArgs,
		},
		{
			name:    "negative_lastn",
			input:   "olist --clientID=12345 --lastN=-1",
			mockFn:  func(f fields) {},
			wantRes: "",
			wantErr: errParseArgs,
		},
		{
			name:  "success",
			input: "olist --clientID=12345 --lastN=1 --inStorage",
			mockFn: func(f fields) {
				f.orderRepo.EXPECT().GetOrdersByFilter(gomock.Any(), order.Filter{
					ClientID:     12345,
					PerPage:      1,
					PageN:        0,
					SortedByDate: true,
					Statuses:     []order.Status{order.Received, order.Refunded},
				}).Return([]order.Order{{OrderID: 1, ClientID: 12345, Status: order.Received}}, nil)
			},
			wantRes: "\nOrders list:\n  orderID=1 clientID=12345 status=received\n",
			wantErr: nil,
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

			res, err := controller.ClientOrders(ctx, strings.Split(tt.input, " ")[1:])

			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.wantRes, res)
		})
	}
}
