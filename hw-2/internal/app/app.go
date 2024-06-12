package app

import (
	"context"
	"os/signal"
	"syscall"

	"gitlab.ozon.dev/antonkraeww/homeworks/hw-2/internal/app/config"
	"gitlab.ozon.dev/antonkraeww/homeworks/hw-2/internal/repository"
	"gitlab.ozon.dev/antonkraeww/homeworks/hw-2/internal/service"
	hashgen "gitlab.ozon.dev/antonkraeww/homeworks/hw-2/internal/workers/hash_generator"
	workerpool "gitlab.ozon.dev/antonkraeww/homeworks/hw-2/internal/workers/worker_pool"
)

func Start(storageFile, configPath string) {
	workersConfig := config.ParseWorkersConfig(configPath)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	orderRepositoryJSON := repository.NewOrderRepository(storageFile)

	hashGen := hashgen.NewHashGenerator(workersConfig.HashesN)
	orderService := service.NewOrderService(orderRepositoryJSON, hashGen)

	workerPool := workerpool.NewWorkerPool(workersConfig.WorkersN, workersConfig.TasksN)
	commands := controllers.NewCLI(orderService, workerPool)

	hashGen.Run(ctx)
	workerPool.Run()
	commands.Run(ctx, cancel)
}
