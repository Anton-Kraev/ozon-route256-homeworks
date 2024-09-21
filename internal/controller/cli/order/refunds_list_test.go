package order

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.ozon.dev/antonkraeww/ozon-route256-homeworks/internal/models/domain/order"
	service "gitlab.ozon.dev/antonkraeww/ozon-route256-homeworks/internal/service/order"
)

func TestOrderController_RefundsList(t *testing.T) {
	t.Parallel()

	type fields struct {
		orderRepo *service.MockorderRepository
		wrapRepo  *service.MockwrapRepository
		hashes    *service.MockhashGenerator
	}

	var (
		ctx      = context.Background()
		testTime = time.Date(
			2006,
			1,
			2,
			15,
			4,
			5,
			0,
			time.FixedZone("", 0),
		)
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
			name:    "negative_pagen",
			input:   "rlist --pageN=-1",
			mockFn:  func(f fields) {},
			wantRes: "",
			wantErr: errParseArgs,
		},
		{
			name:    "negative_perpage",
			input:   "rlist --perPage=-1",
			mockFn:  func(f fields) {},
			wantRes: "",
			wantErr: errParseArgs,
		},
		{
			name:  "success",
			input: "rlist --pageN=1 --perPage=1",
			mockFn: func(f fields) {
				f.orderRepo.EXPECT().GetOrdersByFilter(gomock.Any(), order.Filter{
					PerPage:      1,
					PageN:        1,
					SortedByDate: true,
					Statuses:     []order.Status{order.Refunded},
				}).Return([]order.Order{{
					OrderID:       1,
					ClientID:      12345,
					StatusChanged: testTime,
				}}, nil)
			},
			wantRes: "\nRefunds list:\n  orderID=1 clientID=12345 refunded=2006-01-02 15:04:05 +0000 +0000\n",
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

			res, err := controller.RefundsList(ctx, strings.Split(tt.input, " ")[1:])

			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.wantRes, res)
		})
	}
}
