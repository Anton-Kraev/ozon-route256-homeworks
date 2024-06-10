package main

import "gitlab.ozon.dev/antonkraeww/homeworks/hw-2/internal/app"

func main() {
	app.App{StorageFile: parseArgs()}.Start()
}
