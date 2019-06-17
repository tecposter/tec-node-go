package store

import (
	"testing"
	"bytes"
	"os"
)

func TestSetGet(t *testing.T) {
	t.Log("store Set and Get")

	dir := "./store-test"

	s,err := Open(dir)
	if err != nil {
		t.Error(err)
	}
	defer closeStore(t, s, dir)

	key := "key"
	val := "val"

	s.Set([]byte(key), []byte(val))

	v, err := s.Get([]byte(key))

	if err != nil {
		t.Error(err)
	}

	//v = []byte("fdsa")
	if bytes.Compare([]byte(key), v) == 1 {
		t.Errorf("val expected %s, got %s", val, v)
	}

}

func closeStore(t *testing.T, s *Store, dir string) {
	s.Close()
	err := os.RemoveAll(dir)
	if err != nil {
		t.Error(err)
	}
	//t.Log(v)
}
