package main

import (
	"flag"
	"log"
	"os"
	"strings"
)

func parseArgs() string {
	var (
		args     = os.Args[1:]
		fileName string
	)

	fs := flag.NewFlagSet("app", flag.ContinueOnError)
	fs.StringVar(&fileName, "filename", "orders.json", "use --filename=orders.json")

	if err := fs.Parse(args); err != nil {
		log.Fatalln(err)
	}
	if !strings.HasSuffix(fileName, ".json") {
		log.Fatalln("storage file must be .json")
	}

	return fileName
}
