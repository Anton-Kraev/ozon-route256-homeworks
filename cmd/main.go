package main

import "gitlab.ozon.dev/antonkraeww/homeworks/internal/app"

func main() {
	app.App{StorageFile: parseArgs()}.Start()
}
