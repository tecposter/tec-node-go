// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"flag"
	"log"
	"net/http"
	"github.com/tecposter/tec-server-go/internal/server/ws"
	"github.com/tecposter/tec-server-go/internal/server/web"
)

type landingHandler struct{}

func (landingHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	web.Home(w, req)
}

var addr = flag.String("addr", "0.0.0.0:8090", "http service address")
func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/ws", ws.Handle)
	http.Handle("/", landingHandler{})
	//http.HandleFunc("/", landing.Home)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

/*
package main

import (
	"fmt"
	"net/http"
)

type apiHandler struct{}

func (apiHandler) ServeHTTP(http.ResponseWriter, *http.Request) {}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/api/", apiHandler{})
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		// The "/" pattern matches everything, so we need to check
		// that we're at the root here.
		if req.URL.Path != "/" {
			http.NotFound(w, req)
			return
		}
		fmt.Fprintf(w, "Welcome to the home page!")
	})
}

*/
