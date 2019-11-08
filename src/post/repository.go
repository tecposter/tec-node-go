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

func (repo *repository) saveContent(id dto.ID, typeID int, content string) error {
	has, err := repo.hasContentID(id)
	if err != nil {
		return err
	}
	if has {
		return nil
	}

	stmt, err := repo.db.Prepare("insert into content(id, typeID, content) values (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id, typeID, content)
	if err != nil {
		return err
	}

	return nil
}

func (repo *repository) hasDraftByPostID(postID dto.ID) (bool, error) {
	stmt, err := repo.db.Prepare("select id from draft where id = ? limit 1")
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(postID)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	return rows.Next(), nil
}

func (repo *repository) hasContentID(contentID dto.ID) (bool, error) {
	stmt, err := repo.db.Prepare("select id from content where id = ? limit 1")
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(contentID)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	return rows.Next(), nil
}

func (repo *repository) lastCommit(postID dto.ID) (*commitDTO, error) {
	stmt, err := repo.db.Prepare("select id, postID, contentID, created from [commit] where postID = ? order by created desc limit 1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var c commitDTO
	rows, err := stmt.Query(postID)
	if err != nil {
		return &c, err
	}
	defer rows.Close()

	if !rows.Next() {
		return &c, nil
	}

	err = rows.Scan(&c.ID, &c.PostID, &c.ContentID, &c.Created)
	return &c, err
}

func (repo *repository) commit(commitID dto.ID, postID dto.ID, contentID dto.ID) error {
	last, err := repo.lastCommit(postID)
	if err != nil {
		return err
	}

	if contentID.Equal(last.ContentID) {
		return errContentNotChange
	}

	has, err := repo.hasDraftByPostID(postID)
	if err != nil {
		return err
	}
	if !has {
		return errDraftNotFound
	}

	tx, err := repo.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Commit()

	now := time.Now().UnixNano()
	stmt1, err := tx.Prepare("insert into [commit](id, postID, contentID, created) values (?, ?, ?, ?)")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt1.Close()
	_, err = stmt1.Exec(commitID, postID, contentID, now)
	if err != nil {
		tx.Rollback()
		return err
	}

	stmt2, err := tx.Prepare("update post set commitID = ? where id = ?")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt2.Close()
	_, err = stmt2.Exec(commitID, postID)
	if err != nil {
		tx.Rollback()
		return err
	}

	stmt3, err := tx.Prepare("delete from draft where id = ?")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt3.Close()
	_, err = stmt3.Exec(postID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// err = tx.Commit()
	return nil
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
