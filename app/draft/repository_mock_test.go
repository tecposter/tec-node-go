package draft

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
		mock.ExpectPrepare("insert into draft").
			WillReturnError(errors.New("some error"))

		repo := newRepo(db)
		err := repo.insert(newDraft(dto.ID("id"), "content"))
		assert.NotNil(t, err)
	})
}

func TestMockUpdate(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	var sqlPattern = "update draft set changed (.+), content"

	t.Run("Should return error when executing db.Prepare failed", func(t *testing.T) {
		expectedErr := errors.New("db.Prepare failed")
		mock.ExpectPrepare(sqlPattern).
			WillReturnError(expectedErr)

		repo := newRepo(db)
		err := repo.update(newDraft(dto.ID("id"), "content"))
		assert.Equal(t, expectedErr, err)
	})

	t.Run("Should return error when executing stmt.Exec failed", func(t *testing.T) {
		expectedErr := errors.New("stmt.Exec failed")
		d := newDraft(dto.ID("id-1234"), "content-11234")
		mock.ExpectPrepare(sqlPattern).
			ExpectExec().
			WithArgs(d.Changed, d.Content, d.ID).
			WillReturnError(expectedErr)

		repo := newRepo(db)
		err := repo.update(d)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("Should return error when getting ErrorResult from stmt.Exec", func(t *testing.T) {
		expectedErr := errors.New("stmt.Exec return ErrorResult")
		d := newDraft(dto.ID("id-1234"), "content-11234")
		mock.ExpectPrepare(sqlPattern).
			ExpectExec().
			WithArgs(d.Changed, d.Content, d.ID).
			WillReturnResult(sqlmock.NewErrorResult(expectedErr))

		repo := newRepo(db)
		err := repo.update(d)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("Should return error when affecting no rows", func(t *testing.T) {
		d := newDraft(dto.ID("id-1234"), "content-11234")
		mock.ExpectPrepare(sqlPattern).
			ExpectExec().
			WithArgs(d.Changed, d.Content, d.ID).
			WillReturnResult(sqlmock.NewResult(0, 0))

		repo := newRepo(db)
		err := repo.update(d)
		assert.Equal(t, errAffectNoRows, err)
	})
}

func TestMockHas(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	var sqlPattern = "select id from draft where id = (.+) limit 1"

	t.Run("Should return error when executing db.Prepare failed", func(t *testing.T) {
		expectedErr := errors.New("db.Prepare failed")
		mock.ExpectPrepare(sqlPattern).
			WillReturnError(expectedErr)

		repo := newRepo(db)
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

		repo := newRepo(db)
		_, err := repo.has(id)
		assert.Equal(t, expectedErr, err)
	})
}

func TestMockFecth(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	var sqlPattern = "select id, changed, content from draft where id = (.+) limit 1"

	t.Run("Should return error when executing db.Prepare failed", func(t *testing.T) {
		expectedErr := errors.New("db.Prepare failed")
		mock.ExpectPrepare(sqlPattern).
			WillReturnError(expectedErr)

		repo := newRepo(db)
		_, err := repo.fetch(dto.ID("id"))
		assert.Equal(t, expectedErr, err)
	})
}
