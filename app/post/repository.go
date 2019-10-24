package post

import (
	"database/sql"
	"time"

	"github.com/tecposter/tec-node-go/lib/dto"
)

type repository struct {
	db *sql.DB
}

func newRepo(db *sql.DB) *repository {
	return &repository{db: db}
}

func (repo *repository) create(postID dto.ID) error {
	tx, err := repo.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Commit()

	now := time.Now().UnixNano()
	stmt1, err := tx.Prepare("insert into post(id, commitID, created) values (?, ?, ?)")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt1.Close()
	_, err = stmt1.Exec(postID, dto.EmptyID(), now)
	if err != nil {
		tx.Rollback()
		return err
	}

	stmt2, err := tx.Prepare("insert into draft(id, content, changed) values (?, ?, ?)")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt2.Close()
	_, err = stmt2.Exec(postID, "", now)
	if err != nil {
		tx.Rollback()
		return err
	}

	// err = tx.Commit()
	return nil
}

/*
tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("insert into foo(id, name) values(?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	for i := 0; i < 100; i++ {
		_, err = stmt.Exec(i, fmt.Sprintf("こんにちわ世界%03d", i))
		if err != nil {
			log.Fatal(err)
		}
	}
	tx.Commit()
*/
