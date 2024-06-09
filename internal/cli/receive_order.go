package cli

import (
	"errors"
	"flag"
	"gitlab.ozon.dev/antonkraeww/homeworks/internal/models/requests"
	"time"
)

func (c *CLI) receiveOrder(args []string) (string, error) {
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

	storedUntil, errTime := time.Parse(timeFormat, storedUntilStr)
	if errTime != nil {
		return "", errTime
	}

	return "", c.Service.ReceiveOrder(requests.ReceiveOrderRequest{
		OrderID:     orderID,
		ClientID:    clientID,
		StoredUntil: storedUntil,
	})
}
