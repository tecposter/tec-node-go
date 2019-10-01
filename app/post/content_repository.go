package post

import (
	"database/sql"
	"errors"
	"github.com/tecposter/tec-node-go/lib/dto"
)

var errAffectNoRows = errors.New("Affect No Rows")

type contentRepository struct {
	db *sql.DB
}

func newContentRepo(db *sql.DB) *contentRepository {
	return &contentRepository{db: db}
}

func (repo *contentRepository) insert(c *contentDTO) error {
	stmt, err := repo.db.Prepare("insert into content(id, type, created, content) values (?, ?, ?, ?)")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(c.ID, c.Type, c.Created, c.Content)
	return err
}

/*
func (repo *contentRepository) update(c *content) error {
	stmt, err := repo.db.Prepare("update content set type = ?, created = ?, content = ? where id = ?")
	if err != nil {
		return err
	}
	res, err := stmt.Exec(c.Type, c.Created, c.Content, c.ID)
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
*/

func (repo *contentRepository) has(id dto.ID) (bool, error) {
	stmt, err := repo.db.Prepare("select id from content where id = ? limit 1")
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

func (repo *contentRepository) fetch(id dto.ID) (*contentDTO, error) {
	stmt, err := repo.db.Prepare("select id, type, created, content from content where id = ? limit 1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var c contentDTO
	err = stmt.QueryRow(id).Scan(&c.ID, &c.Type, &c.Created, &c.Content)
	return &c, err
}
