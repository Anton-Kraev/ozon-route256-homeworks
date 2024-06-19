package app

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"gitlab.ozon.dev/antonkraeww/homeworks/hw-3/internal/app/config"
	controller "gitlab.ozon.dev/antonkraeww/homeworks/hw-3/internal/controllers/order"
	"gitlab.ozon.dev/antonkraeww/homeworks/hw-3/internal/pg"
	repository "gitlab.ozon.dev/antonkraeww/homeworks/hw-3/internal/repository/order"
	service "gitlab.ozon.dev/antonkraeww/homeworks/hw-3/internal/service/order"
	hashgen "gitlab.ozon.dev/antonkraeww/homeworks/hw-3/internal/workers/hash_generator"
	workerpool "gitlab.ozon.dev/antonkraeww/homeworks/hw-3/internal/workers/worker_pool"
	"gitlab.ozon.dev/antonkraeww/homeworks/hw-3/middlewares"
)

func Start(configPath string) {
	workersConfig := config.ParseWorkersConfig(configPath)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	connPool, err := pg.NewPoolConn(ctx)
	if err != nil {
		log.Fatalln("can't open postgres connection pool")
	}
	defer connPool.Close()

	var (
		orderRepository = repository.NewOrderRepository(connPool)
		hashGen         = hashgen.NewHashGenerator(workersConfig.HashesN)
		orderService    = service.NewOrderService(orderRepository, hashGen)
		workerPool      = workerpool.NewWorkerPool(workersConfig.WorkersN, workersConfig.TasksN)
		txMiddleware    = middlewares.NewTransactionMiddleware(connPool)
		commands        = controller.NewCLI(orderService, workerPool, txMiddleware)
	)

	hashGen.Run(ctx)
	workerPool.Run()
	commands.Run(ctx, cancel)
}
