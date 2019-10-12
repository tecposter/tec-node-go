package app

import (
	//	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/tecposter/tec-node-go/lib/db/sqlite3"
)

// Application in app
type Application struct {
	// db       *sql.DB
	dataDir  string
	bindAddr string
}

var upgrader = websocket.Upgrader{}

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
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade: ", err)
			return
		}
		defer c.Close()

		db, err := sqlite3.Open(app.dataDir)
		if err != nil {
			log.Print("sqlite3.Open: ", err)
			return
		}
		defer db.Close()

		for {
			mt, message, err := c.ReadMessage()
			if err != nil {
				log.Print("read: ", err)
				continue
			}

			log.Printf("recv: %s", message)
			r, err := unmarshalWSReq(message)
			if err != nil {
				log.Println("unmarshalWSReq: ", err)
				continue
			}

			b, err := r.Marshal()
			if err != nil {
				log.Println("r.Marshal: ", err)
				continue
			}

			err = c.WriteMessage(mt, b)
			if err != nil {
				log.Println("write: ", err)
				continue
			}
		}
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