package main

import (
	"errors"
	"flag"
	"os"
	"strings"
)

func parseArgs() (string, error) {
	var (
		args     = os.Args[1:]
		fileName string
	)

	fs := flag.NewFlagSet("app", flag.ContinueOnError)
	fs.StringVar(&fileName, "filename", "orders.json", "use --filename=orders.json")

	if err := fs.Parse(args); err != nil {
		return "", err
	}
	if !strings.HasSuffix(fileName, ".json") {
		return "", errors.New("storage file must be .json")
	}

	return fileName, nil
}
