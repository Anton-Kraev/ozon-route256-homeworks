package cli

import (
	"errors"
	"flag"
	"time"
)

func (c CLI) receiveOrder(args []string) error {
	const timeFormat = "02.01.2006-15:04:05"
	var (
		orderID, clientID uint64
		storedUntilStr    string
	)

	fs := flag.NewFlagSet(receiveOrder, flag.ContinueOnError)
	fs.Uint64Var(&orderID, "orderID", 0, "use --orderID=12345")
	fs.Uint64Var(&clientID, "clientID", 0, "use --clientID=67890")
	fs.StringVar(&storedUntilStr, "storedUntil", "", "use --storedUntil="+timeFormat)

	if err := fs.Parse(args); err != nil {
		return err
	}
	if orderID == 0 {
		return errors.New("orderID must be positive number")
	}
	if clientID == 0 {
		return errors.New("clientID must be positive number")
	}

	storedUntil, errTime := time.Parse(timeFormat, storedUntilStr)
	if errTime != nil {
		return errTime
	}

	return c.Module.ReceiveOrder(orderID, clientID, storedUntil)
}
