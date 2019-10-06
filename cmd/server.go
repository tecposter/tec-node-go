package main

import (
	"flag"
	"github.com/tecposter/tec-node-go/lib/iotool"
	"log"
	"path"
)

const (
	defaultBindAddr = "127.0.0.1:7890"
)

func main() {
	currDir, err := iotool.CurrDir()
	if err != nil {
		log.Fatal(err)
	}
	defaultDataDir := path.Join(currDir, "data")
	dataDir := flag.String("datadir", defaultDataDir, "Data Directory")
	bindAddr := flag.String("bind", defaultBindAddr, "Bind Addr")
	flag.Parse()

	log.Printf("data directory: %s, binding addr: %s", *dataDir, *bindAddr)
}
