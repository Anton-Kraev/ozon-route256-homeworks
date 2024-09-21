package order

import (
	"context"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	errsdomain "github.com/Anton-Kraev/ozon-route256-homeworks/internal/models/domain/errors"
	service "github.com/Anton-Kraev/ozon-route256-homeworks/internal/service/order"
)

func TestOrderController_ReceiveOrder(t *testing.T) {
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
			name:    "no_orderid",
			input:   "receive --clientID=1 --weight=1 --cost=1 --storedUntil=02.01.2006-15:04:05",
			mockFn:  func(f fields) {},
			wantErr: errParseArgs,
		},
		{
			name:    "no_clientid",
			input:   "receive --orderID=1 --weight=1 --cost=1 --storedUntil=02.01.2006-15:04:05",
			mockFn:  func(f fields) {},
			wantErr: errParseArgs,
		},
		{
			name:    "no_weight",
			input:   "receive --orderID=1 --clientID=1 --cost=1 --storedUntil=02.01.2006-15:04:05",
			mockFn:  func(f fields) {},
			wantErr: errParseArgs,
		},
		{
			name:    "no_cost",
			input:   "receive --orderID=1 --clientID=1 --weight=1 --storedUntil=02.01.2006-15:04:05",
			mockFn:  func(f fields) {},
			wantErr: errParseArgs,
		},
		{
			name:    "bad_time_format",
			input:   "receive --orderID=1 --clientID=1 --weight=1 --cost=1 --storedUntil=02:01:2006_15.04.03",
			mockFn:  func(f fields) {},
			wantErr: errParseArgs,
		},
		{
			name:    "successful_parse",
			input:   "receive --orderID=1 --clientID=1 --weight=1 --cost=1 --wrap=box --storedUntil=02.01.2006-15:04:05",
			mockFn:  func(f fields) {},
			wantErr: errsdomain.ErrRetentionTimeInPast,
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

			res, err := controller.ReceiveOrder(ctx, strings.Split(tt.input, " ")[1:])

			require.ErrorIs(t, err, tt.wantErr)
			assert.Empty(t, res)
		})
	}
}
