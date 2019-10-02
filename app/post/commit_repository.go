package post

import (
	"database/sql"
	"github.com/tecposter/tec-node-go/lib/dto"
)

type commitRepository struct {
	db *sql.DB
}

func newCommitRepo(db *sql.DB) *commitRepository {
	return &commitRepository{db: db}
}

func (repo *commitRepository) insert(c *commitDTO) error {
	stmt, err := repo.db.Prepare("insert into [commit](id, postID, contentID, created) values (?, ?, ?, ?)")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(c.ID, c.PostID, c.ContentID, c.Created)
	return err
}

func (repo *commitRepository) has(id dto.ID) (bool, error) {
	stmt, err := repo.db.Prepare("select id from [commit] where id = ? limit 1")
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

func (repo *commitRepository) fetch(id dto.ID) (*commitDTO, error) {
	stmt, err := repo.db.Prepare("select id, postID, contentID, created from [commit] where id = ? limit 1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var c commitDTO
	err = stmt.QueryRow(id).Scan(&c.ID, &c.PostID, &c.ContentID, &c.Created)
	return &c, err
}
