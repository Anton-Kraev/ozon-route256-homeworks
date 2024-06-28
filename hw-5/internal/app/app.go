package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"

	"gitlab.ozon.dev/antonkraeww/homeworks/hw-5/internal/app/config"
	"gitlab.ozon.dev/antonkraeww/homeworks/hw-5/internal/cli"
	orderCtrl "gitlab.ozon.dev/antonkraeww/homeworks/hw-5/internal/controller/cli/order"
	wrapCtrl "gitlab.ozon.dev/antonkraeww/homeworks/hw-5/internal/controller/cli/wrap"
	"gitlab.ozon.dev/antonkraeww/homeworks/hw-5/internal/middlewares"
	"gitlab.ozon.dev/antonkraeww/homeworks/hw-5/internal/pg"
	orderRepo "gitlab.ozon.dev/antonkraeww/homeworks/hw-5/internal/repository/order"
	wrapRepo "gitlab.ozon.dev/antonkraeww/homeworks/hw-5/internal/repository/wrap"
	orderSrvc "gitlab.ozon.dev/antonkraeww/homeworks/hw-5/internal/service/order"
	wrapSrvc "gitlab.ozon.dev/antonkraeww/homeworks/hw-5/internal/service/wrap"
	hashgen "gitlab.ozon.dev/antonkraeww/homeworks/hw-5/internal/workers/hash_generator"
	workerpool "gitlab.ozon.dev/antonkraeww/homeworks/hw-5/internal/workers/worker_pool"
)

const (
	envFilePath       = "../.env"
	workersConfigPath = "../configs/workers.json"
)

func Start() {
	if err := godotenv.Load(envFilePath); err != nil {
		log.Fatalln("can't load environment file")
	}

	workersConfig := config.ParseWorkersConfig(workersConfigPath)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	connPool, err := pg.NewPoolConn(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalln("can't open postgres connection pool")
	}
	defer connPool.Close()

	var (
		orderRepository = orderRepo.NewOrderRepository(connPool)
		wrapRepository  = wrapRepo.NewWrapRepository(connPool)

		hashGen      = hashgen.NewHashGenerator(workersConfig.HashesN)
		orderService = orderSrvc.NewOrderService(orderRepository, wrapRepository, hashGen)
		wrapService  = wrapSrvc.NewWrapService(wrapRepository)

		orderController = orderCtrl.NewOrderController(orderService)
		wrapController  = wrapCtrl.NewWrapController(wrapService)

		workerPool   = workerpool.NewWorkerPool(workersConfig.WorkersN, workersConfig.TasksN)
		txMiddleware = middlewares.NewTransactionMiddleware(connPool)
		commands     = cli.NewCLI(orderController, wrapController, workerPool, txMiddleware)
	)

	hashGen.Run(ctx)
	workerPool.Run()
	commands.Run(ctx, cancel)
}
