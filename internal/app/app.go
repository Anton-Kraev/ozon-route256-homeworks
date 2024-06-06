package app

import (
	"gitlab.ozon.dev/antonkraeww/homeworks/internal/controller/cli"
	"gitlab.ozon.dev/antonkraeww/homeworks/internal/repository"
	"gitlab.ozon.dev/antonkraeww/homeworks/internal/service"
)

type App struct {
	StorageFile string
}

func (app App) Start() {
	repositoryJSON := repository.NewOrderRepository(app.StorageFile)
	deliveryPointService := service.NewOrderService(repositoryJSON)
	commands := cli.NewCLI(&deliveryPointService)
	commands.Run()
}
