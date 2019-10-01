package post

import (
	"github.com/stretchr/testify/assert"
	"github.com/tecposter/tec-node-go/lib/db/sqlite3"
	"github.com/tecposter/tec-node-go/lib/dto"
	"os"
	"testing"
)

const dirMode = 0755

var dbDir = "./test-db-post-content-dir"

func TestContentRepo(t *testing.T) {
	err := os.Mkdir(dbDir, dirMode)
	assert.Nil(t, err)
	defer os.RemoveAll(dbDir)

	db, err := sqlite3.Open(dbDir)
	assert.Nil(t, err)
	defer db.Close()

	id1 := dto.ID("id1-fdafda")
	content1 := "content1-some random"
	type1 := typeMarkdown
	c1 := newContent(id1, type1, content1)

	t.Run("Should insert successfully", func(t *testing.T) {
		repo := newContentRepo(db)
		err = repo.insert(c1)
		assert.Nil(t, err)
	})

	t.Run("Should has specified id", func(t *testing.T) {
		repo := newContentRepo(db)
		has, err := repo.has(c1.ID)
		assert.Nil(t, err)
		assert.Equal(t, true, has, "Should has id: [%s]", c1.ID)
	})

	t.Run("Fetched item should equal with sepcified item", func(t *testing.T) {
		repo := newContentRepo(db)
		fetchedC1, err := repo.fetch(c1.ID)
		assert.Nil(t, err)
		assert.Equal(t, c1.Created, fetchedC1.Created, "Created should equal")
		assert.Equal(t, c1.Type, fetchedC1.Type, "Type should equal")
		assert.Equal(t, c1.Content, fetchedC1.Content, "Content should equal")
	})

	/*
		t.Run("Should update successfully", func(t *testing.T) {
			repo := newContentRepo(db)
			newContent := "new-content: hdfiewohfdsa"
			newC := newContent(id1, typeHTML, newContent)
			err := repo.update(newC)
			assert.Nil(t, err)

			fetched, err := repo.fetch(id1)
			assert.Nil(t, err)
			assert.Equal(t, newContent, fetched.Content, "Content should equal")
			assert.Greater(t, fetched.Created, c1.Created, "[Created] should increased")
		})
	*/

	t.Run("Should return error when inserting with duplicated id", func(t *testing.T) {
		c := newContent(id1, typeHTML, "any-ffgdf")
		repo := newContentRepo(db)
		err := repo.insert(c)
		assert.NotNil(t, err)
	})

	/*
		t.Run("Should return error when updating with not-existed id", func(t *testing.T) {
			d := newDraft(dto.ID("random-234fdsafad"), "any-content-21432")
			repo := newContentRepo(db)
			err := repo.update(d)
			assert.NotNil(t, err)
		})
	*/
}
