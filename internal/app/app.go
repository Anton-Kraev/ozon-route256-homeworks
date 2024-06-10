package app

import (
	"context"
	"os/signal"
	"syscall"

	"gitlab.ozon.dev/antonkraeww/homeworks/internal/app/config"
	"gitlab.ozon.dev/antonkraeww/homeworks/internal/cli"
	"gitlab.ozon.dev/antonkraeww/homeworks/internal/repository"
	"gitlab.ozon.dev/antonkraeww/homeworks/internal/service"
	hashgen "gitlab.ozon.dev/antonkraeww/homeworks/internal/workers/hash_generator"
	workerpool "gitlab.ozon.dev/antonkraeww/homeworks/internal/workers/worker_pool"
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
