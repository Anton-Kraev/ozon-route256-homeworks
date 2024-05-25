package main

import (
	"fmt"
	"log"

	"gitlab.ozon.dev/antonkraeww/homeworks/hw-1/internal/cli"
	"gitlab.ozon.dev/antonkraeww/homeworks/hw-1/internal/module"
	"gitlab.ozon.dev/antonkraeww/homeworks/hw-1/internal/storage"
)

const (
	fileName = "orders.json"
)

func main() {
	storageJSON := storage.NewStorage(fileName)
	deliveryPointService := module.NewModule(module.Deps{
		Storage: storageJSON,
	})
	commands := cli.NewCLI(cli.Deps{
		Module: deliveryPointService,
	})

	if err := commands.Run(); err != nil {
		log.Fatalln(err)
	}
	fmt.Println("done")
}
