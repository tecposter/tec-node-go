package post

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/tecposter/tec-node-go/lib/dto"
	"testing"
)

func TestMockInsert(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	t.Run("Should return error when executing db.Prepare failed", func(t *testing.T) {
		mock.ExpectPrepare("insert into post").
			WillReturnError(errors.New("some error"))

		repo := newContentRepo(db)
		err := repo.insert(newContent(dto.ID("id"), typeHTML, "content"))
		assert.NotNil(t, err)
	})
}

func TestMockHas(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	var sqlPattern = "select id from content where id = (.+) limit 1"

	t.Run("Should return error when executing db.Prepare failed", func(t *testing.T) {
		expectedErr := errors.New("db.Prepare failed")
		mock.ExpectPrepare(sqlPattern).
			WillReturnError(expectedErr)

		repo := newContentRepo(db)
		_, err := repo.has(dto.ID("id"))
		assert.Equal(t, expectedErr, err)
	})

	t.Run("Should return error when executing stmt.Query failed", func(t *testing.T) {
		expectedErr := errors.New("stmt.Query failed")
		id := dto.ID("id-1234")
		mock.ExpectPrepare(sqlPattern).
			ExpectQuery().
			WithArgs(id).
			WillReturnError(expectedErr)

		repo := newContentRepo(db)
		_, err := repo.has(id)
		assert.Equal(t, expectedErr, err)
	})
}

func TestMockFecth(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	var sqlPattern = "select id, type, created, content from content where id = (.+) limit 1"

	t.Run("Should return error when executing db.Prepare failed", func(t *testing.T) {
		expectedErr := errors.New("db.Prepare failed")
		mock.ExpectPrepare(sqlPattern).
			WillReturnError(expectedErr)

		repo := newContentRepo(db)
		_, err := repo.fetch(dto.ID("id"))
		assert.Equal(t, expectedErr, err)
	})
}
