package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/Anton-Kraev/ozon-route256-homeworks/internal/app/config"
	"github.com/Anton-Kraev/ozon-route256-homeworks/internal/cli"
	orderCtrl "github.com/Anton-Kraev/ozon-route256-homeworks/internal/controller/cli/order"
	wrapCtrl "github.com/Anton-Kraev/ozon-route256-homeworks/internal/controller/cli/wrap"
	"github.com/Anton-Kraev/ozon-route256-homeworks/internal/middlewares"
	"github.com/Anton-Kraev/ozon-route256-homeworks/internal/pg"
	eventRepo "github.com/Anton-Kraev/ozon-route256-homeworks/internal/repository/event"
	orderRepo "github.com/Anton-Kraev/ozon-route256-homeworks/internal/repository/order"
	wrapRepo "github.com/Anton-Kraev/ozon-route256-homeworks/internal/repository/wrap"
	orderSrvc "github.com/Anton-Kraev/ozon-route256-homeworks/internal/service/order"
	wrapSrvc "github.com/Anton-Kraev/ozon-route256-homeworks/internal/service/wrap"
	hashgen "github.com/Anton-Kraev/ozon-route256-homeworks/internal/workers/hash_generator"
	outbox "github.com/Anton-Kraev/ozon-route256-homeworks/internal/workers/outbox_processor"
	workerpool "github.com/Anton-Kraev/ozon-route256-homeworks/internal/workers/worker_pool"
)

const (
	outboxHandlePeriod = 3 * time.Second

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
		eventRepository = eventRepo.NewEventRepository(connPool)

		hashGen      = hashgen.NewHashGenerator(workersConfig.HashesN)
		orderService = orderSrvc.NewOrderService(orderRepository, wrapRepository, hashGen)
		wrapService  = wrapSrvc.NewWrapService(wrapRepository)

		orderController = orderCtrl.NewOrderController(orderService)
		wrapController  = wrapCtrl.NewWrapController(wrapService)

		workerPool   = workerpool.NewWorkerPool(workersConfig.WorkersN, workersConfig.TasksN)
		txMiddleware = middlewares.NewTransactionMiddleware(connPool)
		commands     = cli.NewCLI(eventRepository, orderController, wrapController, workerPool, txMiddleware)

		outboxProcessor = outbox.NewOutboxProcessor(eventRepository)
	)

	outboxProcessor.Start(ctx, outboxHandlePeriod)
	hashGen.Run(ctx)
	workerPool.Run()
	commands.Run(ctx, cancel)
}
