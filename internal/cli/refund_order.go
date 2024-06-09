package cli

import (
	"errors"
	"flag"
	"gitlab.ozon.dev/antonkraeww/homeworks/internal/models/requests"
)

func (c *CLI) refundOrder(args []string) error {
	var orderID, clientID uint64

	fs := flag.NewFlagSet(refundOrder, flag.ContinueOnError)
	fs.Uint64Var(&orderID, "orderID", 0, "use --orderID=12345")
	fs.Uint64Var(&clientID, "clientID", 0, "use --clientID=67890")

	if err := fs.Parse(args); err != nil {
		return err
	}
	if orderID == 0 {
		return errors.New("orderID must be positive number")
	}
	if clientID == 0 {
		return errors.New("clientID must be positive number")
	}

	req := requests.RefundOrderRequest{OrderID: orderID, ClientID: clientID}
	c.cmdManager.AddTask(func() (string, error) {
		return refundOrderTask(c, req)
	})

	return nil
}

func refundOrderTask(cli *CLI, req requests.RefundOrderRequest) (string, error) {
	hash, ok := cli.cmdManager.GetHash()
	if !ok {
		return "", errors.New("hash generation stopped")
	}

	req.Hash = hash

	cli.cmdManager.Mutex.Lock()
	err := cli.Service.RefundOrder(req)
	cli.cmdManager.Mutex.Unlock()

	return "", err
}
