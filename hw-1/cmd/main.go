package main

import (
	"log"

	"gitlab.ozon.dev/antonkraeww/homeworks/hw-1/internal/cli"
	"gitlab.ozon.dev/antonkraeww/homeworks/hw-1/internal/module"
	"gitlab.ozon.dev/antonkraeww/homeworks/hw-1/internal/storage"
)

func main() {
	fileName, err := parseArgs()
	if err != nil {
		log.Fatalln(err)
	}

	storageJSON := storage.NewOrderStorage(fileName)
	deliveryPointService := module.NewOrderModule(storageJSON)
	commands := cli.NewCLI(&deliveryPointService)

	commands.Run()
}
