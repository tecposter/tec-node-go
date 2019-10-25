package draft

import (
	"database/sql"
	"errors"
	"time"

	"github.com/tecposter/tec-node-go/lib/dto"
)

var (
	errPostIDNotExists = errors.New("Post ID not exists")
	errAffectNoRows    = errors.New("Affect No Rows")
)

type repository struct {
	db *sql.DB
}

func newRepo(db *sql.DB) *repository {
	return &repository{db: db}
}

func (repo *repository) save(postID dto.ID, content string) error {
	ok, err := repo.has(postID)
	if err != nil {
		return err
	}
	if !ok {
		return errPostIDNotExists
	}

	stmt, err := repo.db.Prepare("update draft set changed = ?, content = ? where id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	changed := time.Now().UnixNano()
	res, err := stmt.Exec(changed, content, postID)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return errAffectNoRows
	}
	return nil
}

func (repo *repository) has(id dto.ID) (bool, error) {
	stmt, err := repo.db.Prepare("select id from draft where id = ? limit 1")
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(id)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	return rows.Next(), nil
}

func (repo *repository) fetch(id dto.ID) (*draftDTO, error) {
	stmt, err := repo.db.Prepare("select id, changed, content from draft where id = ? limit 1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var d draftDTO
	err = stmt.QueryRow(id).Scan(&d.ID, &d.Changed, &d.Content)
	return &d, err
}
