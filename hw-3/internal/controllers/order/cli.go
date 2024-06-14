package order

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"gitlab.ozon.dev/antonkraeww/homeworks/hw-2/internal/models/domain/order"
)

const timeFormat = "02.01.2006-15:04:05"

type orderService interface {
	ClientOrders(clientID uint64, lastN uint, inStorage bool) ([]order.Order, error)
	DeliverOrders(ordersID []uint64) error
	ReceiveOrder(orderID, clientID uint64, storedUntil time.Time) error
	RefundsList(pageN, perPage uint) ([]order.Order, error)
	RefundOrder(orderID, clientID uint64) error
	ReturnOrder(orderID uint64) error
}

type workerPool interface {
	AddTask(taskID int, command string, task func() (string, error))
	SetNumWorkers(workersN int)
	Shutdown()
	Done() <-chan struct{}
}

type CLI struct {
	Service           orderService
	availableCommands []command
	cmdCounter        int
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
func (c *CLI) Run(ctx context.Context, cancel context.CancelFunc) {
	fmt.Println("The application is running")
	fmt.Println("Type help to get a list of available commands")

	scanner := bufio.NewScanner(os.Stdin)

	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Stopping application...")
				c.workerPool.Shutdown()

				return
			default:
				c.handleCommand(scanner, cancel)
			}
		}
	}()

	// waiting for all tasks processing
	<-c.workerPool.Done()
	fmt.Println("The application has been stopped")
}

func (c *CLI) handleCommand(scanner *bufio.Scanner, cancel context.CancelFunc) {
	for scanner.Scan() {
		input := strings.Split(strings.TrimSpace(scanner.Text()), " ")
		if len(input) == 0 || input[0] == "" {
			continue
		}

		if input[0] == exit {
			cancel()
		} else {
			c.runCommand(input)
		}

		break
	}
}

func (c *CLI) runCommand(input []string) {
	var (
		inputString = strings.Join(input, " ")
		comm        = input[0]
		args        = input[1:]
	)

	switch comm {
	case help:
		c.help()
	case numWorkers:
		c.numWorkers(args)
	case receiveOrder:
		c.cmdCounter++
		c.workerPool.AddTask(c.cmdCounter, inputString, func() (string, error) {
			c.mutex.Lock()
			defer c.mutex.Unlock()

			return c.receiveOrder(args)
		})
	case returnOrder:
		c.cmdCounter++
		c.workerPool.AddTask(c.cmdCounter, inputString, func() (string, error) {
			c.mutex.Lock()
			defer c.mutex.Unlock()

			return c.returnOrder(args)
		})
	case deliverOrders:
		c.cmdCounter++
		c.workerPool.AddTask(c.cmdCounter, inputString, func() (string, error) {
			c.mutex.Lock()
			defer c.mutex.Unlock()

			return c.deliverOrders(args)
		})
	case clientOrders:
		c.cmdCounter++
		c.workerPool.AddTask(c.cmdCounter, inputString, func() (string, error) {
			return c.clientOrders(args)
		})
	case refundOrder:
		c.cmdCounter++
		c.workerPool.AddTask(c.cmdCounter, inputString, func() (string, error) {
			c.mutex.Lock()
			defer c.mutex.Unlock()

			return c.refundOrder(args)
		})
	case refundsList:
		c.cmdCounter++
		c.workerPool.AddTask(c.cmdCounter, inputString, func() (string, error) {
			return c.refundsList(args)
		})
	default:
		fmt.Printf("unknown command %s\n", comm)
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
	fs.UintVar(&workersN, "num", 0, "use --num=4")

	if err := fs.Parse(args); err != nil {
		fmt.Println(err)
	}
	if workersN == 0 {
		fmt.Println("number of workers must be more than zero")
	}

	c.workerPool.SetNumWorkers(int(workersN))
}
