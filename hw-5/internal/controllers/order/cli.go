package order

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"

	"gitlab.ozon.dev/antonkraeww/homeworks/hw-4/internal/models/domain/order"
)

const (
	requestTimeout = 5 * time.Second
	timeFormat     = "02.01.2006-15:04:05"
)

type orderService interface {
	ClientOrders(ctx context.Context, clientID uint64, lastN uint, inStorage bool) ([]order.Order, error)
	DeliverOrders(ctx context.Context, ordersID []uint64) error
	ReceiveOrder(ctx context.Context, order order.Order) error
	RefundsList(ctx context.Context, pageN, perPage uint) ([]order.Order, error)
	RefundOrder(ctx context.Context, orderID, clientID uint64) error
	ReturnOrder(ctx context.Context, orderID uint64) error
}

type workerPool interface {
	AddTask(taskID int, command string, task func() (string, error))
	SetNumWorkers(workersN int)
	Shutdown()
	Done() <-chan struct{}
}

type txMiddleware interface {
	CreateTransactionContext(
		ctx context.Context,
		txOptions pgx.TxOptions,
		args []string,
		handler func(ctx context.Context, args []string) (string, error),
	) (res string, err error)
}

type CLI struct {
	Service           orderService
	availableCommands []command
	cmdCounter        int
	workerPool        workerPool
	txMiddleware      txMiddleware
}

func NewCLI(service orderService, pool workerPool, middleware txMiddleware) CLI {
	return CLI{
		Service:           service,
		availableCommands: commandsList,
		workerPool:        pool,
		txMiddleware:      middleware,
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
				c.handleCommand(cancel, scanner)
			}
		}
	}()

	// waiting for all tasks processing
	<-c.workerPool.Done()
	fmt.Println("The application has been stopped")
}

func (c *CLI) handleCommand(cancel context.CancelFunc, scanner *bufio.Scanner) {
	for scanner.Scan() {
		input := strings.Split(strings.TrimSpace(scanner.Text()), " ")
		if len(input) == 0 || input[0] == "" {
			continue
		}

		switch input[0] {
		case exit:
			cancel()
		case help:
			c.help()
		case numWorkers:
			c.numWorkers(input[1:])
		default:
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

		txReadOptions  = pgx.TxOptions{AccessMode: pgx.ReadOnly, IsoLevel: pgx.ReadCommitted}
		txWriteOptions = pgx.TxOptions{AccessMode: pgx.ReadWrite, IsoLevel: pgx.RepeatableRead}

		parentCtx   = context.Background()
		ctx, cancel = context.WithTimeout(parentCtx, requestTimeout)
	)

	switch comm {
	case receiveOrder:
		c.cmdCounter++
		c.workerPool.AddTask(c.cmdCounter, inputString, func() (string, error) {
			defer cancel()

			return c.txMiddleware.CreateTransactionContext(ctx, txWriteOptions, args, c.receiveOrder)
		})
	case returnOrder:
		c.cmdCounter++
		c.workerPool.AddTask(c.cmdCounter, inputString, func() (string, error) {
			defer cancel()

			return c.txMiddleware.CreateTransactionContext(ctx, txWriteOptions, args, c.returnOrder)
		})
	case deliverOrders:
		c.cmdCounter++
		c.workerPool.AddTask(c.cmdCounter, inputString, func() (string, error) {
			defer cancel()

			return c.txMiddleware.CreateTransactionContext(ctx, txWriteOptions, args, c.deliverOrders)
		})
	case clientOrders:
		c.cmdCounter++
		c.workerPool.AddTask(c.cmdCounter, inputString, func() (string, error) {
			defer cancel()

			return c.txMiddleware.CreateTransactionContext(ctx, txReadOptions, args, c.clientOrders)
		})
	case refundOrder:
		c.cmdCounter++
		c.workerPool.AddTask(c.cmdCounter, inputString, func() (string, error) {
			defer cancel()

			return c.txMiddleware.CreateTransactionContext(ctx, txWriteOptions, args, c.refundOrder)
		})
	case refundsList:
		c.cmdCounter++
		c.workerPool.AddTask(c.cmdCounter, inputString, func() (string, error) {
			defer cancel()

			return c.txMiddleware.CreateTransactionContext(ctx, txReadOptions, args, c.refundsList)
		})
	default:
		fmt.Printf("unknown command %s\n", comm)
		cancel()
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
