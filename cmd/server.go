package main

import (
	"flag"

	"github.com/tecposter/tec-node-go/src/app"
)

const (
	defaultConfigJSONFile = "~/.tec/config.json"
)

func main() {
	f := flag.String("config", defaultConfigJSONFile, "Config JSON File")
	flag.Parse()

	app.Run(*f)
}
