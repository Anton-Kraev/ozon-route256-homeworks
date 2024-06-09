package cli

import (
	"bufio"
	"context"
	"fmt"
	"gitlab.ozon.dev/antonkraeww/homeworks/internal/models/domain/order"
	"gitlab.ozon.dev/antonkraeww/homeworks/internal/models/requests"
	"gitlab.ozon.dev/antonkraeww/homeworks/internal/models/tasks"
	"os"
	"strings"
)

const timeFormat = "02.01.2006-15:04:05"

type orderService interface {
	ReceiveOrder(req requests.ReceiveOrderRequest) error
	ReturnOrder(req requests.ReturnOrderRequest) error
	DeliverOrders(req requests.DeliverOrdersRequest) error
	ClientOrders(req requests.ClientOrdersRequest) ([]order.Order, error)
	RefundOrder(req requests.RefundOrderRequest) error
	RefundsList(req requests.RefundsListRequest) ([]order.Order, error)
}

type taskLogger interface {
	Log(taskRes tasks.TaskResult)
}

type workerPool interface {
	AddTask(taskID int, task func() (string, error))
	GetTaskResult() tasks.TaskResult
}

type CLI struct {
	Service           orderService
	availableCommands []command
	logger            taskLogger
	workerPool        workerPool
}

func NewCLI(service orderService, logger taskLogger, pool workerPool) CLI {
	return CLI{
		Service:           service,
		availableCommands: commandsList,
		logger:            logger,
		workerPool:        pool,
	}
}

// Run runs command-line application, processes entered commands.
func (c *CLI) Run(ctx context.Context) {
	fmt.Println("The application is running")
	fmt.Println("Type help to get a list of available commands")
	defer fmt.Println("The application has been stopped")

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		comm := strings.Split(strings.TrimSpace(scanner.Text()), " ")
		if len(comm) == 0 || comm[0] == "" {
			continue
		}
		if comm[0] == exit {
			break
		}

		c.handleCommand(comm[0], comm[1:])
	}
}

func (c *CLI) handleCommand(comm string, args []string) {
	switch comm {
	case help:
		c.help()
	case receiveOrder:
		c.workerPool.AddTask(1, func() (string, error) {
			return c.receiveOrder(args)
		})
	case returnOrder:
		c.workerPool.AddTask(1, func() (string, error) {
			return c.returnOrder(args)
		})
	case deliverOrders:
		c.workerPool.AddTask(1, func() (string, error) {
			return c.deliverOrders(args)
		})
	case clientOrders:
		c.workerPool.AddTask(1, func() (string, error) {
			return c.clientOrders(args)
		})
	case refundOrder:
		c.workerPool.AddTask(1, func() (string, error) {
			return c.refundOrder(args)
		})
	case refundsList:
		c.workerPool.AddTask(1, func() (string, error) {
			return c.refundsList(args)
		})
	default:
		c.unknownCommand(comm)
	}
}

func (c *CLI) help() {
	fmt.Println("\nAvailable commands list:")

	for _, cmd := range c.availableCommands {
		fmt.Printf("  name: %s\n", cmd.name)
		fmt.Printf("  description: %s\n", cmd.desc)
		fmt.Printf("  usage: %s\n\n", cmd.usage)
	}
}

func (c *CLI) unknownCommand(command string) {
	fmt.Printf("unknown command %s\n", command)
}
