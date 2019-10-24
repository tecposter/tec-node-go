package post

import (
	"database/sql"
	"github.com/tecposter/tec-node-go/lib/dto"
)

type postRepository struct {
	db *sql.DB
}

func newPostRepo(db *sql.DB) *postRepository {
	return &postRepository{db: db}
}

func (repo *postRepository) insert(p *postDTO) error {
	stmt, err := repo.db.Prepare("insert into post(id, commitID, created) values (?, ?, ?)")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(p.ID, p.CommitID, p.Created)
	return err
}

func (repo *postRepository) has(id dto.ID) (bool, error) {
	stmt, err := repo.db.Prepare("select id from post where id = ? limit 1")
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

func (repo *postRepository) fetch(id dto.ID) (*postDTO, error) {
	stmt, err := repo.db.Prepare("select id, commitID, created from post where id = ? limit 1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var p postDTO
	err = stmt.QueryRow(id).Scan(&p.ID, &p.CommitID, &p.Created)
	return &p, err
}

func (repo *postRepository) update(p *postDTO) error {
	stmt, err := repo.db.Prepare("update post set commitID = ? where id = ?")
	if err != nil {
		return err
	}

	res, err := stmt.Exec(p.CommitID, p.ID)
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
