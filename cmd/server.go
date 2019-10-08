package main

import (
	"flag"
	"github.com/tecposter/tec-node-go/lib/iotool"
	"log"
	"net/http"
	"path"
)

const (
	defaultBindAddr = "127.0.0.1:7890"
)

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "static/home.html")
}

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

	http.HandleFunc("/", serveHome)

	err = http.ListenAndServe(*bindAddr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
