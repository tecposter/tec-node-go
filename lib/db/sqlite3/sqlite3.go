package sqlite3

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3" // SQLite3 implementation of SQL
	"github.com/tecposter/tec-node-go/lib/iotool"
	"path"
)

var dbFile = "tec.db"
var driverName = "sqlite3"
var sqlStmt = `
create table post (
	id BLOB not null primary key,
	commitID BLOB not null
);
create table [commit] (
	id BLOB not null primary key,
	postID BLOB not null,
	contentID BLOB not null,
	created NUMERIC not null
);
create table content (
	id BLOB not null primary key,
	content TEXT
);
create table draft (
	postID BLOB not null primary key,
	content TEXT,
	changed NUMERIC not null
);

create index idx_commit_postID on [commit] (postID);
`

// SQLite3 Errors
var (
	ErrDBDIRNotFound = errors.New("DB dir not found")
)

// Open returns sql.DB for sqlite3
func Open(dir string) (*sql.DB, error) {
	if !iotool.FileExists(dir) {
		return nil, ErrDBDIRNotFound
	}

	dbPath := path.Join(dir, dbFile)
	existed := iotool.FileExists(dbPath)

	db, err := sql.Open(driverName, dbPath)
	if err != nil {
		return nil, err
	}

	if !existed {
		db.Exec(sqlStmt)
	}
	return db, nil
}
