package cli

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"gitlab.ozon.dev/antonkraeww/homeworks/internal/models/domain/order"
	"gitlab.ozon.dev/antonkraeww/homeworks/internal/models/requests"
	"os"
	"strings"
	"sync"
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

type workerPool interface {
	AddTask(taskID int, task func() (string, error))
	SetNumWorkers(workersN int)
	Done() <-chan struct{}
}

type CLI struct {
	Service           orderService
	availableCommands []command
	workerPool        workerPool
	mutex             sync.Mutex
}

func NewCLI(service orderService, pool workerPool) CLI {
	return CLI{
		Service:           service,
		availableCommands: commandsList,
		workerPool:        pool,
	}
}

// Run runs command-line application, processes entered commands.
func (c *CLI) Run(cancel context.CancelFunc) {
	fmt.Println("The application is running")
	fmt.Println("Type help to get a list of available commands")

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		comm := strings.Split(strings.TrimSpace(scanner.Text()), " ")
		if len(comm) == 0 || comm[0] == "" {
			continue
		}
		if comm[0] == exit {
			cancel()
			break
		}

		c.handleCommand(comm[0], comm[1:])
	}

	<-c.workerPool.Done()
	fmt.Println("The application has been stopped")
}

func (c *CLI) handleCommand(comm string, args []string) {
	switch comm {
	case help:
		c.help()
	case numWorkers:
		c.numWorkers(args)
	case receiveOrder:
		c.workerPool.AddTask(1, func() (string, error) {
			c.mutex.Lock()
			defer c.mutex.Unlock()
			return c.receiveOrder(args)
		})
	case returnOrder:
		c.workerPool.AddTask(1, func() (string, error) {
			c.mutex.Lock()
			defer c.mutex.Unlock()
			return c.returnOrder(args)
		})
	case deliverOrders:
		c.workerPool.AddTask(1, func() (string, error) {
			c.mutex.Lock()
			defer c.mutex.Unlock()
			return c.deliverOrders(args)
		})
	case clientOrders:
		c.workerPool.AddTask(1, func() (string, error) {
			return c.clientOrders(args)
		})
	case refundOrder:
		c.workerPool.AddTask(1, func() (string, error) {
			c.mutex.Lock()
			defer c.mutex.Unlock()
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

func (c *CLI) numWorkers(args []string) {
	var workersN uint

	fs := flag.NewFlagSet(numWorkers, flag.ContinueOnError)
	fs.UintVar(&workersN, "workersN", 0, "use --num=4")

	if err := fs.Parse(args); err != nil {
		fmt.Println(err)
	}
	if workersN == 0 {
		fmt.Println("number of workers must be more than zero")
	}

	c.workerPool.SetNumWorkers(int(workersN))
}

func (c *CLI) unknownCommand(command string) {
	fmt.Printf("unknown command %s\n", command)
}
