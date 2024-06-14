package app

import (
	"context"
	"os/signal"
	"syscall"

	"gitlab.ozon.dev/antonkraeww/homeworks/hw-3/internal/app/config"
	controller "gitlab.ozon.dev/antonkraeww/homeworks/hw-3/internal/controllers/order"
	repository "gitlab.ozon.dev/antonkraeww/homeworks/hw-3/internal/repository/order"
	service "gitlab.ozon.dev/antonkraeww/homeworks/hw-3/internal/service/order"
	hashgen "gitlab.ozon.dev/antonkraeww/homeworks/hw-3/internal/workers/hash_generator"
	workerpool "gitlab.ozon.dev/antonkraeww/homeworks/hw-3/internal/workers/worker_pool"
)

func Start(storageFile, configPath string) {
	workersConfig := config.ParseWorkersConfig(configPath)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	orderRepositoryJSON := repository.NewOrderRepository(storageFile)

	hashGen := hashgen.NewHashGenerator(workersConfig.HashesN)
	orderService := service.NewOrderService(orderRepositoryJSON, hashGen)

	workerPool := workerpool.NewWorkerPool(workersConfig.WorkersN, workersConfig.TasksN)
	commands := controller.NewCLI(orderService, workerPool)

	hashGen.Run(ctx)
	workerPool.Run()
	commands.Run(ctx, cancel)
}
