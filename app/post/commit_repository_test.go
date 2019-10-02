package post

import (
	"github.com/stretchr/testify/assert"
	"github.com/tecposter/tec-node-go/lib/db/sqlite3"
	"github.com/tecposter/tec-node-go/lib/dto"
	"os"
	"testing"
)

func TestCommitRepo(t *testing.T) {
	const dirMode = 0755
	var dbDir = "./test-db-post-commit-dir"

	err := os.Mkdir(dbDir, dirMode)
	assert.Nil(t, err)
	defer os.RemoveAll(dbDir)

	db, err := sqlite3.Open(dbDir)
	assert.Nil(t, err)
	defer db.Close()

	id1 := dto.ID("id1-fdafda")
	postID1 := dto.ID("post-id1-fadf")
	contentID1 := dto.ID("content-id1-fadfda")
	c1 := newCommit(id1, postID1, contentID1)

	t.Run("Should insert successfully", func(t *testing.T) {
		repo := newCommitRepo(db)
		err = repo.insert(c1)
		assert.Nil(t, err)
	})

	t.Run("Should has specified id", func(t *testing.T) {
		repo := newCommitRepo(db)
		has, err := repo.has(c1.ID)
		assert.Nil(t, err)
		assert.Equal(t, true, has, "Should has id: [%s]", c1.ID)
	})

	t.Run("Fetched item should equal with sepcified item", func(t *testing.T) {
		repo := newCommitRepo(db)
		fetchedC1, err := repo.fetch(c1.ID)
		assert.Nil(t, err)
		assert.Equal(t, c1.Created, fetchedC1.Created, "Created should equal")
		assert.Equal(t, c1.PostID, fetchedC1.PostID, "PostID should equal")
		assert.Equal(t, c1.ContentID, fetchedC1.ContentID, "ContentID should equal")
		assert.Equal(t, c1.Created, fetchedC1.Created, "Created should equal")
	})

	t.Run("Should return error when inserting with duplicated id", func(t *testing.T) {
		c := newCommit(id1, dto.ID("post-id-any"), dto.ID("content-id-any"))
		repo := newCommitRepo(db)
		err := repo.insert(c)
		assert.NotNil(t, err)
	})
}
