package main

import (
	"gitlab.ozon.dev/antonkraeww/homeworks/internal/cli"
	"gitlab.ozon.dev/antonkraeww/homeworks/internal/module"
	"gitlab.ozon.dev/antonkraeww/homeworks/internal/storage"
)

func main() {
	fileName := parseArgs()

	storageJSON := storage.NewOrderStorage(fileName)
	deliveryPointService := module.NewOrderModule(storageJSON)
	commands := cli.NewCLI(&deliveryPointService)

	commands.Run()
}
