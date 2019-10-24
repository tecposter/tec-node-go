package post

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/tecposter/tec-node-go/lib/dto"
	"testing"
)

func TestMockPostInsert(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	t.Run("Should return error when executing db.Prepare failed", func(t *testing.T) {
		mock.ExpectPrepare("insert into post").
			WillReturnError(errors.New("some error"))

		repo := newPostRepo(db)
		err := repo.insert(newPost(dto.ID("id"), dto.ID("commit-id")))
		assert.NotNil(t, err)
	})
}

func TestMockPostHas(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	var sqlPattern = "select id from post where id = (.+) limit 1"

	t.Run("Should return error when executing db.Prepare failed", func(t *testing.T) {
		expectedErr := errors.New("db.Prepare failed")
		mock.ExpectPrepare(sqlPattern).
			WillReturnError(expectedErr)

		repo := newPostRepo(db)
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

		repo := newPostRepo(db)
		_, err := repo.has(id)
		assert.Equal(t, expectedErr, err)
	})
}

func TestMockPostFecth(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	var sqlPattern = "select id, commitID, created from post where id = (.+) limit 1"

	t.Run("Should return error when executing db.Prepare failed", func(t *testing.T) {
		expectedErr := errors.New("db.Prepare failed")
		mock.ExpectPrepare(sqlPattern).
			WillReturnError(expectedErr)

		repo := newPostRepo(db)
		_, err := repo.fetch(dto.ID("id"))
		assert.Equal(t, expectedErr, err)
	})
}

func TestMockPostUpdate(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	var sqlPattern = "update post set commitID = (.+) where id = (.+)"

	t.Run("Should return error when executing db.Prepare failed", func(t *testing.T) {
		expectedErr := errors.New("db.Prepare failed")
		mock.ExpectPrepare(sqlPattern).
			WillReturnError(expectedErr)

		repo := newPostRepo(db)
		err := repo.update(newPost(dto.ID("id"), dto.ID("commit-id")))
		assert.Equal(t, expectedErr, err)
	})

	t.Run("Should return error when executing stmt.Exec failed", func(t *testing.T) {
		expectedErr := errors.New("stmt.Exec failed")
		p := newPost(dto.ID("id-1234"), dto.ID("commit-id-11234"))
		mock.ExpectPrepare(sqlPattern).
			ExpectExec().
			WithArgs(p.CommitID, p.ID).
			WillReturnError(expectedErr)

		repo := newPostRepo(db)
		err := repo.update(p)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("Should return error when getting ErrorResult from stmt.Exec", func(t *testing.T) {
		expectedErr := errors.New("stmt.Exec return ErrorResult")
		p := newPost(dto.ID("id-1234"), dto.ID("commit-id-11234"))
		mock.ExpectPrepare(sqlPattern).
			ExpectExec().
			WithArgs(p.CommitID, p.ID).
			WillReturnResult(sqlmock.NewErrorResult(expectedErr))

		repo := newPostRepo(db)
		err := repo.update(p)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("Should return error when affecting no rows", func(t *testing.T) {
		p := newPost(dto.ID("id-1234"), dto.ID("commit-id-11234"))
		mock.ExpectPrepare(sqlPattern).
			ExpectExec().
			WithArgs(p.CommitID, p.ID).
			WillReturnResult(sqlmock.NewResult(0, 0))

		repo := newPostRepo(db)
		err := repo.update(p)
		assert.Equal(t, errAffectNoRows, err)
	})
}
