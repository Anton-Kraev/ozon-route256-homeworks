package cli

import (
	"errors"
	"flag"
	"gitlab.ozon.dev/antonkraeww/homeworks/internal/models/requests"
)

func (c *CLI) returnOrder(args []string) error {
	var orderID uint64

	fs := flag.NewFlagSet(returnOrder, flag.ContinueOnError)
	fs.Uint64Var(&orderID, "orderID", 0, "use --orderID=12345")

	if err := fs.Parse(args); err != nil {
		return err
	}
	if orderID == 0 {
		return errors.New("orderID must be positive number")
	}

	req := requests.ReturnOrderRequest{OrderID: orderID}
	c.cmdManager.AddTask(func() (string, error) {
		return returnOrderTask(c, req)
	})

	return nil
}

func returnOrderTask(cli *CLI, req requests.ReturnOrderRequest) (string, error) {
	hash, ok := cli.cmdManager.GetHash()
	if !ok {
		return "", errors.New("hash generation stopped")
	}

	req.Hash = hash

	cli.cmdManager.Mutex.Lock()
	err := cli.Service.ReturnOrder(req)
	cli.cmdManager.Mutex.Unlock()

	return "", err
}
