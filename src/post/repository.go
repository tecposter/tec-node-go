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
	stmt, err := repo.db.Prepare("select id, postID, contentID, committed from [commit] where postID = ? order by committed desc limit 1")
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

	err = rows.Scan(&c.ID, &c.PostID, &c.ContentID, &c.Committed)
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
	stmt1, err := tx.Prepare("insert into [commit](id, postID, contentID, committed) values (?, ?, ?, ?)")
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

func (repo *repository) create(postID dto.ID) (*postDTO, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Commit()

	now := time.Now().UnixNano()
	stmt1, err := tx.Prepare("insert into post(id, commitID, posted) values (?, ?, ?)")
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	defer stmt1.Close()
	_, err = stmt1.Exec(postID, dto.EmptyID(), now)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	stmt2, err := tx.Prepare("insert into draft(id, content, drafted) values (?, ?, ?)")
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	defer stmt2.Close()
	_, err = stmt2.Exec(postID, "", now)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// err = tx.Commit()
	return &postDTO{
		ID:        postID,
		CommitID:  dto.ID([]byte("")),
		ContentID: dto.ID([]byte("")),
		Content:   "",
		Posted:    now,
		Drafted:   0,
		Committed: 0,
	}, nil
}

func (repo *repository) edit(postID dto.ID) error {
	has, err := repo.hasDraftByPostID(postID)
	if err != nil {
		return err
	}
	if has {
		return nil
	}

	stmt, err := repo.db.Prepare("insert into draft(id, content, drafted) values (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	now := time.Now().UnixNano()
	_, err = stmt.Exec(postID, "", now)
	if err != nil {
		return err
	}

	return nil
}

func (repo *repository) fetch(postID dto.ID) (*postDTO, error) {
	stmt, err := repo.db.Prepare(`select
	p.id, IFNULL(p.commitID, x''), IFNULL(m.contentID, x''), IFNULL(c.content, ''), p.posted, IFNULL(m.committed, 0), IFNULL(d.drafted, 0)
	from post p
	left join [commit] m on m.id = p.commitID
	left join content c on c.id = m.contentID
	left join draft d on d.id = p.id
	where p.id = ? limit 1`)

	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var p postDTO
	err = stmt.QueryRow(postID).Scan(&p.ID, &p.CommitID, &p.ContentID, &p.Content, &p.Posted, &p.Committed, &p.Drafted)
	return &p, err
}

func (repo *repository) list() ([]postItemDTO, error) {
	stmt, err := repo.db.Prepare(`select
	p.id, IFNULL(p.commitID, x''), IFNULL(m.contentID, x''), IFNULL(c.content, ''), IFNULL(p.posted, 0), IFNULL(m.committed, 0), IFNULL(d.drafted, 0)
	from post p
	left join [commit] m on m.id = p.commitID
	left join content c on c.id = m.contentID
	left join draft d on d.id = p.id
	order by m.committed desc`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var arr []postItemDTO
	rows, err := stmt.Query()
	if err != nil {
		return arr, err
	}

	for rows.Next() {
		var p postItemDTO
		var content string
		err = rows.Scan(&p.ID, &p.CommitID, &p.ContentID, &content, &p.Posted, &p.Committed, &p.Drafted)
		if err != nil {
			return arr, err
		}
		p.Title = extractTitle(content)
		arr = append(arr, p)
	}
	return arr, nil
}

func extractTitle(content string) string {
	return content
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
