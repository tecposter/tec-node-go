package draft

import (
	"github.com/stretchr/testify/assert"
	"github.com/tecposter/tec-node-go/lib/db/sqlite3"
	"github.com/tecposter/tec-node-go/lib/dto"
	"os"
	"testing"
)

const dirMode = 0755

var dbDir = "./test-db-fdksaghdfa-dir"

func TestCRUD(t *testing.T) {
	err := os.Mkdir(dbDir, dirMode)
	assert.Nil(t, err)
	defer os.RemoveAll(dbDir)

	db, err := sqlite3.Open(dbDir)
	assert.Nil(t, err)
	defer db.Close()

	id1 := dto.ID("id1-fdafda")
	content1 := "content1-some random"
	d1 := newDraft(id1, content1)

	t.Run("Should insert successfully", func(t *testing.T) {
		repo := newRepo(db)
		err = repo.insert(d1)
		assert.Nil(t, err)
	})

	t.Run("Should has specified id", func(t *testing.T) {
		repo := newRepo(db)
		has, err := repo.has(d1.ID)
		assert.Nil(t, err)
		assert.Equal(t, true, has, "Should has id: [%s]", d1.ID)
	})

	t.Run("Fetched item should equal with sepcified item", func(t *testing.T) {
		repo := newRepo(db)
		fetchedD1, err := repo.fetch(d1.ID)
		assert.Nil(t, err)
		assert.Equal(t, d1.Changed, fetchedD1.Changed, "Changed should equal")
		assert.Equal(t, d1.Content, fetchedD1.Content, "Content should equal")
	})

	t.Run("Should update successfully", func(t *testing.T) {
		repo := newRepo(db)
		newContent := "new-content: hdfiewohfdsa"
		newD := newDraft(id1, newContent)
		err := repo.update(newD)
		assert.Nil(t, err)

		fetched, err := repo.fetch(id1)
		assert.Nil(t, err)
		assert.Equal(t, newContent, fetched.Content, "Content should equal")
		assert.Greater(t, fetched.Changed, d1.Changed, "[Changed] should increased")
	})

	t.Run("Should return error when inserting with duplicated id", func(t *testing.T) {
		d := newDraft(id1, "any-ffgdf")
		repo := newRepo(db)
		err := repo.insert(d)
		assert.NotNil(t, err)
	})

	t.Run("Should return error when updating with not-existed id", func(t *testing.T) {
		d := newDraft(dto.ID("random-234fdsafad"), "any-content-21432")
		repo := newRepo(db)
		err := repo.update(d)
		assert.NotNil(t, err)
	})
}
