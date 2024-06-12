package main

import "gitlab.ozon.dev/antonkraeww/homeworks/hw-2/internal/app"

const configPath = "./configs/"

func main() {
	app.Start(parseArgs(), configPath)
}
