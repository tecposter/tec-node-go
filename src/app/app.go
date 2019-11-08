package app

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"path"

	"github.com/tecposter/tec-node-go/lib/iotool"
	"github.com/tecposter/tec-node-go/src/config"
	"github.com/tecposter/tec-node-go/src/db"
	"github.com/tecposter/tec-node-go/src/draft"
	"github.com/tecposter/tec-node-go/src/post"
	"github.com/tecposter/tec-node-go/src/searcher"
	"github.com/tecposter/tec-node-go/src/ws"
)

const (
	postModule  = "post"
	draftModule = "draft"
)

/*
const (
	cmdUnknown = "cmdUnknown"
)
*/

// errors
var (
	ErrModuleNotFound = errors.New("Module not found")
)

// Run run http server app
func Run(jsonFile string) {
	searcher.Init()
	defer searcher.Close()

	err := config.LoadFromJSONFile(jsonFile)
	if err != nil {
		log.Println("config error: ", err)
	}

	http.HandleFunc("/", handleHome())
	http.HandleFunc("/ws", handleWS())

	err = http.ListenAndServe(config.BindAddr(), nil)
	if err != nil {
		log.Fatal("app.Run: ", err)
	}
}

func handleHome() http.HandlerFunc {
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

func handleWS() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db, err := newDB()
		if err != nil {
			log.Print("newDB: ", err)
		}

		conn, err := ws.NewConn(w, r, db)
		if err != nil {
			log.Print("ws.NewConn: ", err)
		}
		defer conn.Close()

		conn.Handle(func() {
			switch conn.Req().Module() {
			case draftModule:
				draft.Handle(conn)
			case postModule:
				post.Handle(conn)
			default:
				conn.Res().SetErr(ErrModuleNotFound)
			}
		})
	}
}

func newDB() (*sql.DB, error) {
	dataDir := path.Join(config.BaseDir(), "/data")
	iotool.MkdirIfNotExist(dataDir)
	db, err := db.Open(dataDir)
	return db, err
}
