package order

import (
	"context"
	"errors"
	"flag"
	"log"
	"time"

	errsdomain "gitlab.ozon.dev/antonkraeww/homeworks/hw-4/internal/models/domain/errors"
)

func (c *CLI) receiveOrder(ctx context.Context, args []string) (string, error) {
	var (
		orderID, clientID uint64
		storedUntilStr    string
	)

	fs := flag.NewFlagSet(receiveOrder, flag.ContinueOnError)
	fs.Uint64Var(&orderID, "orderID", 0, "use --orderID=12345")
	fs.Uint64Var(&clientID, "clientID", 0, "use --clientID=67890")
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

	storedUntil, err := time.Parse(timeFormat, storedUntilStr)
	if err != nil {
		return "", err
	}

	err = c.Service.ReceiveOrder(ctx, orderID, clientID, storedUntil)
	if err != nil {
		switch {
		case errors.Is(err, errsdomain.ErrOrderIDNotUnique):
			return "", err
		case errors.Is(err, errsdomain.ErrRetentionTimeInPast):
			return "", err
		default:
			log.Println(err.Error())

			return "", errors.New("can't receive order")
		}
	}

	return "", nil
}
