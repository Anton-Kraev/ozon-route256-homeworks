package order

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"

	errsdomain "gitlab.ozon.dev/antonkraeww/ozon-route256-homeworks/internal/models/domain/errors"
)

func (c OrderController) ReturnOrder(ctx context.Context, args []string) (string, error) {
	var orderID uint64

	fs := flag.NewFlagSet("return", flag.ContinueOnError)
	fs.Uint64Var(&orderID, "orderID", 0, "use --orderID=12345")

	if err := fs.Parse(args); err != nil {
		return "", fmt.Errorf("%w: %w", errParseArgs, err)
	}

	if orderID == 0 {
		return "", fmt.Errorf("%w: orderID must be positive number", errParseArgs)
	}

	err := c.service.ReturnOrder(ctx, orderID)
	if err != nil {
		switch {
		case errors.Is(err, errsdomain.ErrOrderNotFound):
			return "", err
		case errors.Is(err, errsdomain.ErrRetentionPeriodNotExpiredYet):
			return "", err
		case errors.Is(err, errsdomain.ErrOrderAlreadyReturned):
			return "", err
		case errors.Is(err, errsdomain.ErrOrderDelivered):
			return "", err
		default:
			log.Println(err.Error())

			return "", errReturnOrder
		}
	}

	return "", nil
}
