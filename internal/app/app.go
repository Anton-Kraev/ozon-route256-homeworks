package app

import (
	"context"
	"gitlab.ozon.dev/antonkraeww/homeworks/internal/cli"
	"gitlab.ozon.dev/antonkraeww/homeworks/internal/repository"
	"gitlab.ozon.dev/antonkraeww/homeworks/internal/service"
	"gitlab.ozon.dev/antonkraeww/homeworks/internal/workers/hash_generator"
	"gitlab.ozon.dev/antonkraeww/homeworks/internal/workers/worker_pool"
	"os/signal"
	"syscall"
)

type App struct {
	StorageFile string
}

func (app App) Start() {
	parentCtx := context.Background()
	ctx, cancel := signal.NotifyContext(parentCtx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	orderRepositoryJSON := repository.NewOrderRepository(app.StorageFile)

	hg := hashgen.NewHashGenerator(5)
	orderService := service.NewOrderService(orderRepositoryJSON, hg)

	wp := workerpool.NewWorkerPool(5, 5)
	commands := cli.NewCLI(orderService, wp)

	hg.Run(ctx)
	wp.Run()
	commands.Run(ctx, cancel)
}
