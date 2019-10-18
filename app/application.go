package app

import (
	//	"database/sql"
	"log"
	"net/http"
)

// Application in app
type Application struct {
	dataDir  string
	bindAddr string
}

// NewApp return Application
func NewApp(dataDir, bindAddr string) *Application {
	return &Application{
		dataDir:  dataDir,
		bindAddr: bindAddr,
	}
}

// Run run app to serve http
func (app *Application) Run() error {
	http.HandleFunc("/", app.handleHome())
	http.HandleFunc("/ws", app.handleWS())

	err := http.ListenAndServe(app.bindAddr, nil)
	return err
}

func (app *Application) handleWS() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := newConn(app.dataDir, w, r)
		if err != nil {
			log.Print("newConn: ", err)
		}
		defer conn.close()

		/*
			conn.onErr(func(err error) {
				log.Print("conn: ", err)
			})
		*/
		conn.run()
	}
}

func (app *Application) handleHome() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
}
