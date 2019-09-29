package content

// https://github.com/DATA-DOG/go-sqlmock/blob/master/README.md

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMockHasID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	t.Run("db.Prepare should fail", func(t *testing.T) {
		expectedErr := errors.New("hasID err")
		mock.ExpectPrepare("select id from content where").
			WillReturnError(expectedErr)

		repo := newRepo(db)
		_, err = repo.hasID([]byte("any"))
		assert.Equal(t, expectedErr, err, "repo.hasID should return error")
	})

	t.Run("stmt.Query should fail", func(t *testing.T) {
		expectedErr := errors.New("stmt.Query return error")
		expectedID := []byte("id")

		mock.ExpectPrepare("select id from content where").
			ExpectQuery().
			WithArgs(expectedID).
			WillReturnError(expectedErr)

		repo := newRepo(db)
		_, err = repo.hasID(expectedID)
		assert.Equal(t, expectedErr, err, "expect stmt.Query failed error")
	})
}

func TestMockFetchContent(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	t.Run("db.prepare should fail", func(t *testing.T) {
		expectedErr := errors.New("db.Prepare error")
		mock.ExpectPrepare("select content from content where").
			WillReturnError(expectedErr)

		repo := newRepo(db)
		_, err = repo.fetchContent([]byte("any"))
		assert.Equal(t, expectedErr, err, "should return error from db.Prepare")
	})
}

func TestMockAdd(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	t.Run("hasID stmt.Query should fail", func(t *testing.T) {
		expectedErr := errors.New("hadID - stmt.Query return error")
		expectedContent := "Test Mock Add"
		expectedID := generateCID(expectedContent)

		mock.ExpectPrepare("select id from content where").
			ExpectQuery().
			WithArgs(expectedID).
			WillReturnError(expectedErr)

		repo := newRepo(db)
		_, err = repo.add(expectedContent)
		assert.Equal(t, expectedErr, err, "expect stmt.Query failed error")
	})

	t.Run("add stmt.Query insert should fail", func(t *testing.T) {
		expectedErr := errors.New("insert error")
		expectedContent := "Test Mock Add"
		expectedID := generateCID(expectedContent)

		mock.ExpectPrepare("select id from content where").
			ExpectQuery().
			WithArgs(expectedID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}))

		mock.ExpectPrepare("insert into content").
			WillReturnError(expectedErr)

		repo := newRepo(db)
		_, err := repo.add(expectedContent)
		assert.Equal(t, expectedErr, err, "expect insert error")
	})
}
