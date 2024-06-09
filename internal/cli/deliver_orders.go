package cli

import (
	"errors"
	"flag"
	"gitlab.ozon.dev/antonkraeww/homeworks/internal/models/requests"
	"strconv"
	"strings"
)

func (c *CLI) deliverOrders(args []string) error {
	var (
		ordersStr  string
		ordersList []uint64
	)

	fs := flag.NewFlagSet(deliverOrders, flag.ContinueOnError)
	fs.StringVar(&ordersStr, "orders", "", "use --orders=1,2,3,4,5")

	if err := fs.Parse(args); err != nil {
		return err
	}

	for _, orderStr := range strings.Split(ordersStr, ",") {
		orderID, err := strconv.ParseUint(orderStr, 10, 64)

		if err != nil {
			return err
		}
		if orderID == 0 {
			return errors.New("orderID must be positive number")
		}

		ordersList = append(ordersList, orderID)
	}

	req := requests.DeliverOrdersRequest{OrdersID: ordersList}
	c.cmdManager.AddTask(func() (string, error) {
		return deliverOrdersTask(c, req)
	})

	return nil
}

func deliverOrdersTask(cli *CLI, req requests.DeliverOrdersRequest) (string, error) {
	req.Hashes = make([]string, len(req.OrdersID))

	for i := range req.OrdersID {
		hash, ok := cli.cmdManager.GetHash()
		if !ok {
			return "", errors.New("hash generation stopped")
		}

		req.Hashes[i] = hash
	}

	cli.cmdManager.Mutex.Lock()
	err := cli.Service.DeliverOrders(req)
	cli.cmdManager.Mutex.Unlock()

	return "", err
}
