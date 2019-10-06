package app

import (
	"database/sql"
	"github.com/tecposter/tec-node-go/lib/db/sqlite3"
)

// Application in app
type Application struct {
	db      *sql.DB
	dataDir string
}

// NewApp return Application
func NewApp(dataDir string) (*Application, error) {
	db, err := sqlite3.Open(dbDir)
	if err != nil {
		return nil, err
	}

	return &Application{
		db:      db,
		dataDir: dataDir}
}

// Close close appliction
func (app *Application) Close() {
	app.db.Close()
}
