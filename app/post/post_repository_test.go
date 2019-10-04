package post

import (
	"github.com/stretchr/testify/assert"
	"github.com/tecposter/tec-node-go/lib/db/sqlite3"
	"github.com/tecposter/tec-node-go/lib/dto"
	"os"
	"testing"
)

func TestPostRepo(t *testing.T) {
	const dirMode = 0755
	var dbDir = "./test-db-post-post-dir"
	err := os.Mkdir(dbDir, dirMode)
	assert.Nil(t, err)
	defer os.RemoveAll(dbDir)

	db, err := sqlite3.Open(dbDir)
	assert.Nil(t, err)
	defer db.Close()

	id1 := dto.ID("id1-fdafda")
	commitID1 := dto.ID("commit-id1-fdsafds")
	p1 := newPost(id1, commitID1)

	t.Run("Should insert successfully", func(t *testing.T) {
		repo := newPostRepo(db)
		err = repo.insert(p1)
		assert.Nil(t, err)
	})

	t.Run("Should has specified id", func(t *testing.T) {
		repo := newPostRepo(db)
		has, err := repo.has(p1.ID)
		assert.Nil(t, err)
		assert.Equal(t, true, has, "Should has id: [%s]", p1.ID)
	})

	t.Run("Fetched item should equal with sepcified item", func(t *testing.T) {
		repo := newPostRepo(db)
		fetchedC1, err := repo.fetch(p1.ID)
		assert.Nil(t, err)
		assert.Equal(t, p1.Created, fetchedC1.Created, "Created should equal")
		assert.Equal(t, p1.CommitID, fetchedC1.CommitID, "CommitID should equal")
	})

	t.Run("Should return error when inserting with duplicated id", func(t *testing.T) {
		p := newPost(id1, dto.ID("commit-id-any-ffgdf"))
		repo := newPostRepo(db)
		err := repo.insert(p)
		assert.NotNil(t, err)
	})

	t.Run("Should update successfully", func(t *testing.T) {
		repo := newPostRepo(db)
		commitID2 := dto.ID("commit-id2: hdfiewohfdsa")
		newP := newPost(id1, commitID2)
		err := repo.update(newP)
		assert.Nil(t, err)

		fetched, err := repo.fetch(id1)
		assert.Nil(t, err)
		assert.Equal(t, commitID2, fetched.CommitID, "CommitID should equal")
	})

	t.Run("Should return error when updating with not-existed id", func(t *testing.T) {
		p := newPost(dto.ID("id-random-234fdsafad"), dto.ID("commit-id-any-21432"))
		repo := newPostRepo(db)
		err := repo.update(p)
		assert.NotNil(t, err)
	})
}
