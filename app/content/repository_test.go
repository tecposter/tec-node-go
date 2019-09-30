package content

import (
	"bytes"
	"github.com/tecposter/tec-node-go/lib/db/sqlite3"
	"os"
	"testing"
)

const dirMode = 0755

var dbDir = "./test-db-dir"

func TestAdd(t *testing.T) {
	err := os.Mkdir(dbDir, dirMode)
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(dbDir)
	db, err := sqlite3.Open(dbDir)
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	t.Run("should return content id", func(t *testing.T) {
		content := "hello world"
		expectedID := generateCID(content)

		repo := newRepo(db)

		id, err := repo.add(content)
		if err != nil {
			t.Error(err)
		}

		if !bytes.Equal(expectedID, id) {
			t.Errorf("id expected: %s, but got %s", expectedID, id)
		}
	})

	t.Run("Same content should return save cid", func(t *testing.T) {
		content := "hello world"
		sameContent := "hello world"
		differentContent := "hello world changed"

		repo := newRepo(db)

		id1, err := repo.add(content)
		if err != nil {
			t.Error(err)
		}
		id2, err := repo.add(sameContent)
		if err != nil {
			t.Error(err)
		}
		id3, err := repo.add(differentContent)
		if err != nil {
			t.Error(err)
		}

		if !bytes.Equal(id1, id2) {
			t.Errorf("Same content should expect same id: %s vs %s", id1, id2)
		}

		if bytes.Equal(id1, id3) {
			t.Errorf("Different content should expect different ids: %s vs %s", id1, id3)
		}
	})

}

func TestFetch(t *testing.T) {
	err := os.Mkdir(dbDir, dirMode)
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(dbDir)
	db, err := sqlite3.Open(dbDir)
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	expectedContent := "hello world"
	expectedID := generateCID(expectedContent)

	repo := newRepo(db)

	id, err := repo.add(expectedContent)
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(expectedID, id) {
		t.Errorf("id expected: %s, but got %s", expectedID, id)
	}

	content, err := repo.fetchContent(id)
	if err != nil {
		t.Error(err)
	}
	if content != expectedContent {
		t.Errorf("content expected: %s, but got %s", expectedContent, content)
	}
}
