package main

import (
	"fmt"
	"os"
	"strings"

	"gitlab.ozon.dev/antonkraeww/homeworks/hw-1/internal/cli"
	"gitlab.ozon.dev/antonkraeww/homeworks/hw-1/internal/module"
	"gitlab.ozon.dev/antonkraeww/homeworks/hw-1/internal/storage"
)

func main() {
	var fileName string
	if len(os.Args) > 1 {
		fileName = os.Args[1]
		if !strings.HasSuffix(fileName, ".json") {
			fmt.Println("storage file must be .json")
			return
		}
	} else {
		fileName = "orders.json"
	}

	storageJSON := storage.NewOrderStorage(fileName)
	deliveryPointService := module.NewOrderModule(storageJSON)
	commands := cli.NewCLI(&deliveryPointService)

	commands.Run()
}
