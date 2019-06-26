package store

import (
	"github.com/dgraph-io/badger"
)

// Iter wraps badger.Iterator
type Iter struct {
	offset int
	limit  int
	db     *DB
}

const (
	limitMax = 9999
)

// SetOffset sets offset
func (wrap *Iter) SetOffset(offset int) *Iter {
	wrap.offset = offset
	return wrap
}

// SetLimit sets limit
func (wrap *Iter) SetLimit(limit int) *Iter {
	wrap.limit = limit
	return wrap
}

// ForEach executes a provided function for each key-value pair
func (wrap *Iter) ForEach(call func([]byte, []byte) error) error {
	return wrap.db.inner.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		offset := wrap.offset
		limit := wrap.limit
		if limit == 0 {
			limit = limitMax
		}

		it.Rewind()
		for offset > 0 && it.Valid() {
			it.Next()
			offset--
		}

		for limit > 0 && it.Valid() {
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

			it.Next()
			limit--
		}

		return nil
	})
}
