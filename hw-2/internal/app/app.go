package app

import (
	"context"
	"os/signal"
	"syscall"

	"gitlab.ozon.dev/antonkraeww/homeworks/hw-2/internal/app/config"
	"gitlab.ozon.dev/antonkraeww/homeworks/hw-2/internal/cli"
	"gitlab.ozon.dev/antonkraeww/homeworks/hw-2/internal/repository"
	"gitlab.ozon.dev/antonkraeww/homeworks/hw-2/internal/service"
	hashgen "gitlab.ozon.dev/antonkraeww/homeworks/hw-2/internal/workers/hash_generator"
	workerpool "gitlab.ozon.dev/antonkraeww/homeworks/hw-2/internal/workers/worker_pool"
)

type App struct {
	StorageFile string
}

func (app App) Start() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	orderRepositoryJSON := repository.NewOrderRepository(app.StorageFile)

	hashGen := hashgen.NewHashGenerator(config.HashesN)
	orderService := service.NewOrderService(orderRepositoryJSON, hashGen)

	workerPool := workerpool.NewWorkerPool(config.WorkersN, config.TasksN)
	commands := cli.NewCLI(orderService, workerPool)

	hashGen.Run(ctx)
	workerPool.Run()
	commands.Run(ctx, cancel)
}
