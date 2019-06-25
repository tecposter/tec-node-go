package store

import (
	"github.com/dgraph-io/badger"
)

type DB struct {
	inner *badger.DB
}

type Txn struct {
	inner *badger.Txn
}

type txnHandler func(*Txn) error

func Open(dirs ...string) (*DB, error) {
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

	return &DB{inner: db}, nil
}

func (db *DB) Close() {
	db.inner.Close()
}

func (db *DB) Get(key []byte) ([]byte, error) {
	var valCopy []byte

	err := db.inner.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
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

func (db *DB) Set(key, val []byte) error {
	err := db.inner.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry(key, val)
		err := txn.SetEntry(e)
		return err
	})

	return err
}

func (db *DB) MultiSet(arr [][2][]byte) error {
	return db.inner.Update(func(txn *badger.Txn) error {
		for _, pair := range arr {
			entry := badger.NewEntry(pair[0], pair[1])
			err := txn.SetEntry(entry)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func Pair(p1, p2 []byte) [2][]byte {
	return [2][]byte{p1, p2}
}

func Arr(arr ...[2][]byte) [][2][]byte {
	return arr
}

func (db *DB) Update(hdl txnHandler) error {
	return db.inner.Update(func(txn *badger.Txn) error {
		return hdl(&Txn{inner: txn})
	})
}

/*
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


func (s *Store) Has(key []byte) bool {
	if _, err := s.Get(key); err == nil {
		return true
	} else if err == badger.ErrKeyNotFound {
		return false
	} else {
		panic(err)
	}
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
*/
