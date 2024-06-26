package order

import (
	"context"
	"errors"
	"flag"
	"log"
	"time"

	errsdomain "gitlab.ozon.dev/antonkraeww/homeworks/hw-5/internal/models/domain/errors"
	"gitlab.ozon.dev/antonkraeww/homeworks/hw-5/internal/models/domain/order"
)

func (c *CLI) receiveOrder(ctx context.Context, args []string) (string, error) {
	var (
		orderID, clientID        uint64
		orderWeight, orderCost   uint
		wrapType, storedUntilStr string
	)

	fs := flag.NewFlagSet(receiveOrder, flag.ContinueOnError)
	fs.Uint64Var(&orderID, "orderID", 0, "use --orderID=12345")
	fs.Uint64Var(&clientID, "clientID", 0, "use --clientID=67890")
	fs.UintVar(&orderWeight, "weight", 0, "use --weight=2000")
	fs.UintVar(&orderCost, "cost", 0, "use --cost=500")
	fs.StringVar(&wrapType, "wrap", "", "use --wrap=100")
	fs.StringVar(&storedUntilStr, "storedUntil", "", "use --storedUntil="+timeFormat)

	if err := fs.Parse(args); err != nil {
		return "", err
	}

	if orderID == 0 {
		return "", errors.New("orderID must be positive number")
	}
	if clientID == 0 {
		return "", errors.New("clientID must be positive number")
	}
	if orderWeight == 0 {
		return "", errors.New("weight must be positive number")
	}
	if orderCost == 0 {
		return "", errors.New("cost must be positive number")
	}

	storedUntil, err := time.Parse(timeFormat, storedUntilStr)
	if err != nil {
		return "", err
	}

	if wrapType == "" {
		wrapType = string(order.Nowrap)
	}

	err = c.Service.ReceiveOrder(ctx, order.Order{
		OrderID:     orderID,
		ClientID:    clientID,
		Weight:      orderWeight,
		Cost:        orderCost,
		WrapType:    order.Wrap(wrapType),
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
		default:
			log.Println(err.Error())

			return "", errors.New("can't receive order")
		}
	}

	return "", nil
}
