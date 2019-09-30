package sqlite3

import (
	"os"
	"testing"
)

func TestOpen(t *testing.T) {
	const dirMode = 0755
	dbDir := "./test-db-dir"
	os.RemoveAll(dbDir)
	err := os.Mkdir(dbDir, dirMode)
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(dbDir)

	t.Run("shoud throw error DB dir not found", func(t *testing.T) {
		errDir := "./err-db-dir"
		_, err = Open(errDir)
		if err != ErrDBDIRNotFound {
			t.Error("Expected error: DB dir not found")
		}
	})

	t.Run("Should create tables: post commit content draft", func(t *testing.T) {
		db, err := Open(dbDir)
		if err != nil {
			t.Error(err)
		}
		defer db.Close()

		rows, err := db.Query(`select name from sqlite_master
		where type = 'table' and name not like 'sqlite_%'`)
		if err != nil {
			t.Error(err)
		}
		defer rows.Close()

		var expectedTables = map[string]struct{}{
			"post":    struct{}{},
			"commit":  struct{}{},
			"content": struct{}{},
			"draft":   struct{}{},
		}
		var currentTables []string

		count := 0
		for rows.Next() {
			var name string
			err = rows.Scan(&name)
			if err != nil {
				t.Error(err)
			}
			if _, ok := expectedTables[name]; !ok {
				t.Errorf("Unexpected table %s", name)
			}
			currentTables = append(currentTables, name)
			count++
		}

		if count != 4 {
			t.Errorf("Table number expected 4, but got %d", count)
		}
	})

	t.Run("Should contain idx_commit_postID", func(t *testing.T) {
		db, err := Open(dbDir)
		if err != nil {
			t.Error(err)
		}
		defer db.Close()

		rows, err := db.Query(`select count(name) from sqlite_master
			where type="index" and name = 'idx_commit_postID'`)

		if err != nil {
			t.Error(err)
		}
		defer rows.Close()

		var count int
		rows.Next()
		err = rows.Scan(&count)

		if err != nil {
			t.Error(err)
		}
		if count != 1 {
			t.Error("Cound not find id_commit_postID")
		}
	})

	t.Run("Should throw unkown driver", func(t *testing.T) {
		driverName = "unkown"
		_, err := Open(dbDir)
		if err == nil {
			t.Error("Expected error unkown driver")
		}
		driverName = "sqlite3"
	})
}
