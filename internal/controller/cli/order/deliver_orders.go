package order

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"

	"github.com/Anton-Kraev/ozon-route256-homeworks/internal/helpers"
	errsdomain "github.com/Anton-Kraev/ozon-route256-homeworks/internal/models/domain/errors"
)

func (c OrderController) DeliverOrders(ctx context.Context, args []string) (string, error) {
	var ordersStr string

	fs := flag.NewFlagSet("deliver", flag.ContinueOnError)
	fs.StringVar(&ordersStr, "orders", "", "use --orders=1,2,3,4,5")

	if err := fs.Parse(args); err != nil {
		return "", fmt.Errorf("%w: %w", errParseArgs, err)
	}

	orders, err := helpers.StrToUint64Arr(ordersStr)
	if err != nil {
		return "", fmt.Errorf("%w: invalid format, must be --orders=<id1>,<id2>,<id3>: %w", errParseArgs, err)
	}

	for _, order := range orders {
		if order == 0 {
			return "", fmt.Errorf("%w: orderID must be positive number", errParseArgs)
		}
	}

	err = c.service.DeliverOrders(ctx, orders)
	if err != nil {
		switch {
		case errors.Is(err, errsdomain.ErrOrderNotFound):
			return "", err
		case errors.Is(err, errsdomain.ErrDifferentClients):
			return "", err
		case errors.Is(err, errsdomain.ErrUnexpectedOrderStatus):
			return "", err
		case errors.Is(err, errsdomain.ErrRetentionPeriodExpired):
			return "", err
		default:
			log.Println(err.Error())

			return "", errDeliverOrders
		}
	}

	return "", nil
}
