package storage

import (
	"github.com/coreos/bbolt"
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

func (bs *boltStore) WriteBatch(batch []Entry) error {
	return bs.db.Update(func(tx *bolt.Tx) error {
		for _, entry := range batch {
			bucket, err := tx.CreateBucketIfNotExists(entry.BucketKey())
			if err != nil {
				return err
			}
			if err := bucket.Put(entry.Key, entry.Value); err != nil {
				return err
			}
		}
		return nil
	})
}

func (bs *boltStore) ReadBatch(batch []*Entry) error {
	return bs.db.View(func(tx *bolt.Tx) error {
		for _, entry := range batch {
			if bucket := tx.Bucket(entry.BucketKey()); bucket != nil {
				entry.Value = bucket.Get(entry.Key)
			}
		}
		return nil
	})
}

func (bs *boltStore) Close() error {
	return bs.db.Close()
}
