package store

import (
	"github.com/dgraph-io/badger"
)

// Iterator wraps badger.Iterator
type Iterator struct {
	opts badger.IteratorOptions
	db   *DB
}

// ForEach executes a provided function for each key-value pair
func (wrap *Iterator) ForEach(call func([]byte, []byte) error) error {
	return wrap.db.inner.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(wrap.opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			var valCopy []byte
			item := it.Item()
			key := item.Key()
			err := item.Value(func(v []byte) error {
				valCopy = append([]byte{}, v...)
				return nil
			})
			if err != nil {
				return err
			}

			err = call(key, valCopy)
			if err != nil {
				return err
			}
		}

		return nil
	})
}
