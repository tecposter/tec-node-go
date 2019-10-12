package main

import (
	"flag"
	"log"
	"path"

	// "github.com/gorilla/websocket"

	"github.com/tecposter/tec-node-go/app"
	"github.com/tecposter/tec-node-go/lib/iotool"
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
	iotool.MkdirIfNotExist(*dataDir)

	err = app.NewApp(*dataDir, *bindAddr).Run()
	if err != nil {
		log.Fatal("App.Run: ", err)
	}
}
