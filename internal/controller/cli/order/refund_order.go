package order

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"

	errsdomain "gitlab.ozon.dev/antonkraeww/ozon-route256-homeworks/internal/models/domain/errors"
)

func (c OrderController) RefundOrder(ctx context.Context, args []string) (string, error) {
	var orderID, clientID uint64

	fs := flag.NewFlagSet("refund", flag.ContinueOnError)
	fs.Uint64Var(&orderID, "orderID", 0, "use --orderID=12345")
	fs.Uint64Var(&clientID, "clientID", 0, "use --clientID=67890")

	if err := fs.Parse(args); err != nil {
		return "", fmt.Errorf("%w: %w", errParseArgs, err)
	}

	if orderID == 0 {
		return "", fmt.Errorf("%w: orderID must be positive number", errParseArgs)
	}
	if clientID == 0 {
		return "", fmt.Errorf("%w: clientID must be positive number", errParseArgs)
	}

	err := c.service.RefundOrder(ctx, orderID, clientID)
	if err != nil {
		switch {
		case errors.Is(err, errsdomain.ErrOrderNotFound):
			return "", err
		case errors.Is(err, errsdomain.ErrOrderAlreadyRefunded):
			return "", err
		case errors.Is(err, errsdomain.ErrOrderNotDeliveredYet):
			return "", err
		case errors.Is(err, errsdomain.ErrOrderDeliveredLongAgo):
			return "", err
		default:
			log.Println(err.Error())

			return "", errRefundOrder
		}
	}

	return "", nil
}
