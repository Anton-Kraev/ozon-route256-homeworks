package cli

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

const requestTimeout = 5 * time.Second

type (
	orderController interface {
		RefundOrder(ctx context.Context, args []string) (string, error)
		RefundsList(ctx context.Context, args []string) (string, error)
		ReturnOrder(ctx context.Context, args []string) (string, error)
		ClientOrders(ctx context.Context, args []string) (string, error)
		ReceiveOrder(ctx context.Context, args []string) (string, error)
		DeliverOrders(ctx context.Context, args []string) (string, error)
	}

	workerPool interface {
		AddTask(taskID int, command string, task func() (string, error))
		SetNumWorkers(workersN int)
		Shutdown()
		Done() <-chan struct{}
	}

	txMiddleware interface {
		CreateTransactionContext(
			ctx context.Context,
			txOptions pgx.TxOptions,
			args []string,
			handler func(ctx context.Context, args []string) (string, error),
		) (res string, err error)
	}
)

type CLI struct {
	cmdCounter   int
	controller   orderController
	workerPool   workerPool
	txMiddleware txMiddleware
}

func NewCLI(controller orderController, pool workerPool, middleware txMiddleware) CLI {
	return CLI{
		controller:   controller,
		workerPool:   pool,
		txMiddleware: middleware,
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

		cmd := commandName(input[0])

		switch {
		case slices.Contains(controlCommandsNames, cmd):
			c.handleControlCommand(cancel, input)
		case slices.Contains(orderCommandsNames, cmd):
			c.handleOrderCommand(input)
		default:
			fmt.Printf("unknown command %s\n", cmd)
		}

		break
	}
}
