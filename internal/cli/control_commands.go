package cli

import (
	"context"
	"flag"
	"fmt"
)

func (c *CLI) handleControlCommand(cancel context.CancelFunc, input []string) {
	switch commandName(input[0]) {
	case exit:
		cancel()
	case help:
		c.help()
	case numWorkers:
		c.numWorkers(input[1:])
	}
}

func (c *CLI) help() {
	fmt.Println("\nAvailable commands list:")

	for _, cmd := range allCommandsInfo {
		fmt.Printf("  name: %s\n", cmd.name)
		fmt.Printf("  description: %s\n", cmd.desc)
		fmt.Printf("  usage: %s\n\n", cmd.usage)
	}
}

func (c *CLI) numWorkers(args []string) {
	var workersN uint

	fs := flag.NewFlagSet(string(numWorkers), flag.ContinueOnError)
	fs.UintVar(&workersN, "num", 0, "use --num=4")

	if err := fs.Parse(args); err != nil {
		fmt.Println(err)
	}
	if workersN == 0 {
		fmt.Println("number of workers must be more than zero")
	}

	c.workerPool.SetNumWorkers(int(workersN))
}
