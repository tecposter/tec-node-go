package content

import (
	"crypto/sha256"
	"database/sql"
	"github.com/tecposter/tec-node-go/dto"
)

type repository struct {
	db *sql.DB
}

func newRepo(db *sql.DB) *repository {
	return &repository{db: db}
}

func (repo *repository) add(content string) (dto.ID, error) {
	id := generateCID(content)
	had, err := repo.hasID(id)
	if err != nil {
		return nil, err
	}

	if had {
		return id, nil
	}

	stmt, err := repo.db.Prepare("insert into content(id, content) values (?, ?)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	_, err = stmt.Exec(id, content)
	return id, err
}

func (repo *repository) hasID(id dto.ID) (bool, error) {
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

func (repo *repository) fetchContent(id dto.ID) (string, error) {
	stmt, err := repo.db.Prepare("select content from content where id = ? limit 1")
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	var content string
	err = stmt.QueryRow(id).Scan(&content)
	return content, err
}

func generateCID(content string) dto.ID {
	h := sha256.New()
	h.Write([]byte(content))
	id := h.Sum(nil)
	return dto.ID(id)
}
