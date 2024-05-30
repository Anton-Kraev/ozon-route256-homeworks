package cli

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"gitlab.ozon.dev/antonkraeww/homeworks/hw-1/internal/models"
)

type Module interface {
	ReceiveOrder(orderID uint64, clientID uint64, storedUntil time.Time) error
	ReturnOrder(orderID uint64) error
	DeliverOrders(ordersID []uint64) error
	ClientOrders(clientID uint64, lastN uint, inStorage bool) ([]models.Order, error)
	RefundOrder(orderID uint64, clientID uint64) error
	RefundsList(pageN uint, perPage uint) ([]models.Order, error)
}

type Deps struct {
	Module Module
}

type CLI struct {
	Deps
	availableCommands []command
}

func NewCLI(d Deps) CLI {
	return CLI{
		Deps:              d,
		availableCommands: commandsList,
	}
}

func (c CLI) Run() {
	fmt.Println("The application is running")
	fmt.Println("Type help to get a list of available commands")
	defer fmt.Println("The application has been stopped")

	scanner := bufio.NewScanner(os.Stdin)
	globalCounter := 1
	fmt.Print("1) ")

	for scanner.Scan() {
		comm := strings.Split(strings.TrimSpace(scanner.Text()), " ")
		if len(comm) == 0 || comm[0] == "" {
			fmt.Printf("%d) ", globalCounter)
			continue
		}
		if comm[0] == exit {
			break
		}

		go func(localCounter int) {
			err := c.handleCommand(comm[0], comm[1:])
			fmt.Println()
			if err != nil {
				fmt.Printf("%d) error: %v\n", localCounter, err)
			} else {
				fmt.Printf("%d) ok\n", localCounter)
			}
			fmt.Printf("%d) ", globalCounter)
		}(globalCounter)

		globalCounter++
	}
}

func (c CLI) handleCommand(comm string, args []string) error {
	switch comm {
	case help:
		c.help()
		return nil
	case receiveOrder:
		return c.receiveOrder(args)
	case returnOrder:
		return c.returnOrder(args)
	case deliverOrders:
		return c.deliverOrders(args)
	case clientOrders:
		return c.clientOrders(args)
	case refundOrder:
		return c.refundOrder(args)
	case refundsList:
		return c.refundsList(args)
	default:
		return fmt.Errorf("unknown command %s", comm)
	}
}

func (c CLI) help() {
	fmt.Println("\nAvailable commands list:")
	for _, cmd := range c.availableCommands {
		fmt.Printf("  name: %s\n", cmd.name)
		fmt.Printf("  description: %s\n", cmd.desc)
		fmt.Printf("  usage: %s\n\n", cmd.usage)
	}
}

func (c CLI) receiveOrder(args []string) error {
	var (
		orderID, clientID uint64
		storedUntilStr    string
	)

	fs := flag.NewFlagSet(receiveOrder, flag.ContinueOnError)
	fs.Uint64Var(&orderID, "orderID", 0, "use --orderID=12345")
	fs.Uint64Var(&clientID, "clientID", 0, "use --clientID=67890")
	fs.StringVar(&storedUntilStr, "storedUntil", "", "use --storedUntil=02.01.2006-15:04:05")

	if err := fs.Parse(args); err != nil {
		return err
	}
	if orderID == 0 {
		return errors.New("orderID must be positive number")
	}
	if clientID == 0 {
		return errors.New("clientID must be positive number")
	}
	storedUntil, errTime := time.Parse("02.01.2006-15:04:05", storedUntilStr)
	if errTime != nil {
		return errTime
	}

	return c.Module.ReceiveOrder(orderID, clientID, storedUntil)
}

func (c CLI) returnOrder(args []string) error {
	var orderID uint64

	fs := flag.NewFlagSet(returnOrder, flag.ContinueOnError)
	fs.Uint64Var(&orderID, "orderID", 0, "use --orderID=12345")

	if err := fs.Parse(args); err != nil {
		return err
	}
	if orderID == 0 {
		return errors.New("orderID must be positive number")
	}

	return c.Module.ReturnOrder(orderID)
}

func (c CLI) deliverOrders(args []string) error {
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

	return c.Module.DeliverOrders(ordersList)
}

func (c CLI) clientOrders(args []string) error {
	var (
		clientID  uint64
		lastN     uint
		inStorage bool
	)

	fs := flag.NewFlagSet(clientOrders, flag.ContinueOnError)
	fs.Uint64Var(&clientID, "clientID", 0, "use --clientID=12345")
	fs.UintVar(&lastN, "lastN", 0, "use --lastN=10")
	fs.BoolVar(&inStorage, "storedUntil", false, "use --inStorage")

	if err := fs.Parse(args); err != nil {
		return err
	}
	if clientID == 0 {
		return errors.New("clientID must be positive number")
	}

	orders, err := c.Module.ClientOrders(clientID, lastN, inStorage)
	if err == nil {
		fmt.Println("\nOrders list:")
		for _, order := range orders {
			fmt.Printf(
				"orderID=%d clientID=%d status=%s\n",
				order.OrderID, order.ClientID, order.Status,
			)
		}
	}

	return err
}

func (c CLI) refundOrder(args []string) error {
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

	return c.Module.RefundOrder(orderID, clientID)
}

func (c CLI) refundsList(args []string) error {
	var pageN, perPage uint

	fs := flag.NewFlagSet(refundsList, flag.ContinueOnError)
	fs.UintVar(&pageN, "pageN", 0, "use --pageN=3")
	fs.UintVar(&perPage, "perPage", 0, "use --perPage=10")

	if err := fs.Parse(args); err != nil {
		return err
	}

	refunds, err := c.Module.RefundsList(pageN, perPage)
	if err == nil {
		fmt.Println("\nRefunds list:")
		for _, refund := range refunds {
			fmt.Printf(
				"orderID=%d clientID=%d refunded=%s\n",
				refund.OrderID, refund.ClientID, refund.StatusChanged,
			)
		}
	}

	return err
}
