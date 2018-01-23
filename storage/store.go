package storage

import (
	"bytes"
	"encoding/binary"
)

const (
	TransactionBucket Bucket = iota + 1
	TransactionMetadataBucket
	MilestoneBucket
	StateDiffBucket
	AddressBucket
	ApproveeBucket
	BundleBucket
	TagBucket
)

var allBuckets = []Bucket{
	TransactionBucket,
	TransactionMetadataBucket,
	MilestoneBucket,
	StateDiffBucket,
	AddressBucket,
	ApproveeBucket,
	BundleBucket,
	TagBucket,
}

var bucketKeys = map[Bucket][]byte{}

func init() {
	// Converty bucket int ids to []byte keys
	for _, bucket := range allBuckets {
		var buf bytes.Buffer
		err := binary.Write(&buf, binary.LittleEndian, int(bucket))
		if err != nil {
			panic(err)
		}
		bucketKeys[bucket] = buf.Bytes()
	}
}

type Bucket int

type Entry struct {
	Bucket Bucket
	Key    []byte
	Value  []byte
}

func (e *Entry) BucketKey() []byte {
	return bucketKeys[e.Bucket]
}

type Store interface {
	// WriteBatch writes a batch of entries to the DB.
	WriteBatch(batch []Entry) error

	// ReadBatch reads a batch of entries from the DB. Upon success, the bytes value for each entry should be set.
	ReadBatch(batch []Entry) error
}

// Write is a convenience method to perform a single entry write.
func Write(s Store, key, value []byte, bucket Bucket) error {
	return s.WriteBatch([]Entry{{Bucket: bucket, Key: key, Value: value}})
}

// Read is a convenience method to perform a single entry read.
func Read(s Store, key []byte, bucket Bucket) ([]byte, error) {
	entry := Entry{Bucket: bucket, Key: key, Value: nil}
	if err := s.ReadBatch([]Entry{entry}); err != nil {
		return nil, err
	}
	return entry.Value, nil
}
