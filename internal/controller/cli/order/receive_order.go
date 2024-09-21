package order

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"time"

	errsdomain "gitlab.ozon.dev/antonkraeww/ozon-route256-homeworks/internal/models/domain/errors"
	"gitlab.ozon.dev/antonkraeww/ozon-route256-homeworks/internal/models/domain/order"
)

func (c OrderController) ReceiveOrder(ctx context.Context, args []string) (string, error) {
	var (
		orderID, clientID        uint64
		orderWeight, orderCost   uint
		wrapType, storedUntilStr string
	)

	fs := flag.NewFlagSet("receive", flag.ContinueOnError)
	fs.Uint64Var(&orderID, "orderID", 0, "use --orderID=12345")
	fs.Uint64Var(&clientID, "clientID", 0, "use --clientID=67890")
	fs.UintVar(&orderWeight, "weight", 0, "use --weight=2000")
	fs.UintVar(&orderCost, "cost", 0, "use --cost=500")
	fs.StringVar(&wrapType, "wrap", "", "use --wrap=100")
	fs.StringVar(&storedUntilStr, "storedUntil", "", "use --storedUntil="+timeFormat)

	if err := fs.Parse(args); err != nil {
		return "", fmt.Errorf("%w: %w", errParseArgs, err)
	}

	if orderID == 0 {
		return "", fmt.Errorf("%w: orderID must be positive number", errParseArgs)
	}
	if clientID == 0 {
		return "", fmt.Errorf("%w: clientID must be positive number", errParseArgs)
	}
	if orderWeight == 0 {
		return "", fmt.Errorf("%w: weight must be positive number", errParseArgs)
	}
	if orderCost == 0 {
		return "", fmt.Errorf("%w: cost must be positive number", errParseArgs)
	}

	storedUntil, err := time.Parse(timeFormat, storedUntilStr)
	if err != nil {
		return "", fmt.Errorf("%w: %w", errParseArgs, err)
	}

	err = c.service.ReceiveOrder(ctx, wrapType, order.Order{
		OrderID:     orderID,
		ClientID:    clientID,
		Weight:      orderWeight,
		Cost:        orderCost,
		StoredUntil: storedUntil,
	})
	if err != nil {
		switch {
		case errors.Is(err, errsdomain.ErrOrderIDNotUnique):
			return "", err
		case errors.Is(err, errsdomain.ErrRetentionTimeInPast):
			return "", err
		case errors.Is(err, errsdomain.ErrUnknownOrderWrapType):
			return "", err
		case errors.Is(err, errsdomain.ErrOrderWeightExceedsLimit):
			return "", err
		case errors.Is(err, errsdomain.ErrWrapNotFound):
			return "", err
		default:
			log.Println(err.Error())

			return "", errReceiveOrder
		}
	}

	return "", nil
}
