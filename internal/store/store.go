package store

import (
	"github.com/dgraph-io/badger"
)

type Store struct {
	db *badger.DB
}

func Open(dirs ...string) (*Store, error) {
	opts := badger.DefaultOptions

	opts.Dir = dirs[0]
	if len(dirs) >= 2 {
		opts.ValueDir = dirs[1]
	} else {
		opts.ValueDir = dirs[0]
	}

  db, err := badger.Open(opts)

	if err != nil {
		return nil, err
	}

	s := Store {db: db}
	return &s, nil
}

func (s *Store) Close() {
	s.db.Close()
}

func (s *Store) Set(key, val []byte) error {
	err := s.db.Update(func (txn *badger.Txn) error {
		e := badger.NewEntry(key, val)
		err := txn.SetEntry(e)
		return err
	})

	return err
}

func (s *Store) Get(key []byte) ([]byte, error) {
	var valCopy []byte

	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil  {
			return err
		}

		err = item.Value(func(val []byte) error {
			valCopy = append([]byte{}, val...)
			return nil
		})

		return err
	})

	return valCopy, err
}
