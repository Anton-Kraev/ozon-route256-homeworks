package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"gitlab.ozon.dev/antonkraeww/homeworks/internal/domain/models"
)

const timeFormat = "02.01.2006-15:04:05"

type orderModule interface {
	ReceiveOrder(orderID uint64, clientID uint64, storedUntil time.Time) error
	ReturnOrder(orderID uint64) error
	DeliverOrders(ordersID []uint64) error
	ClientOrders(clientID uint64, lastN uint, inStorage bool) ([]models.Order, error)
	RefundOrder(orderID uint64, clientID uint64) error
	RefundsList(pageN uint, perPage uint) ([]models.Order, error)
}

type CLI struct {
	Module            orderModule
	availableCommands []command
}

func NewCLI(module orderModule) CLI {
	return CLI{
		Module:            module,
		availableCommands: commandsList,
	}
}

// Run runs command-line application, processes entered commands.
func (c CLI) Run() {
	fmt.Println("The application is running")
	fmt.Println("Type help to get a list of available commands")
	defer fmt.Println("The application has been stopped")

	scanner := bufio.NewScanner(os.Stdin)
	cmdCounter := 1
	fmt.Printf("%d) ", cmdCounter)

	for scanner.Scan() {
		comm := strings.Split(strings.TrimSpace(scanner.Text()), " ")
		if len(comm) == 0 || comm[0] == "" {
			fmt.Printf("%d) ", cmdCounter)
			continue
		}
		if comm[0] == exit {
			break
		}

		err := c.handleCommand(comm[0], comm[1:])
		if err != nil {
			fmt.Printf("%d) error: %v\n", cmdCounter, err)
		} else {
			fmt.Printf("%d) ok\n", cmdCounter)
		}

		cmdCounter++
		fmt.Printf("%d) ", cmdCounter)
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
