package storage

import (
	"github.com/boltdb/bolt"
)

func NewBoltStore(path string) (Store, error) {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}
	return &boltStore{db}, nil
}

type boltStore struct {
	db *bolt.DB
}

func (bs *boltStore) Close() {
	bs.db.Close()
}

func (bs *boltStore) WriteBatch(batch []Entry) error {
	return nil // TODO(era): impl
}

func (bs *boltStore) ReadBatch(batch []Entry) error {
	return nil // TODO(era): impl
}
