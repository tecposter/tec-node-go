package store

import (
	// "errors"
	"database/sql"
	_ "github.com/mattn/go-sqlite3" // sqlite3 database
)

// https://golang.org/pkg/database/sql/#Open

// DB wraps sqlite3
type DB struct{}

// Open returns a new DB object
// func Open(dirs ...string) (*DB, error) {
// }

// Close closes a DB
func (db *DB) Close() {
}

// func Get
// func Set
// func Delete
