package storage

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

type Bucket int

type Entry struct {
	bucket Bucket
	key    []byte
	bytes  []byte
}

type Store interface {
	// WriteBatch writes a batch of entries to the DB.
	WriteBatch(batch []Entry) error

	// ReadBatch reads a batch of entries from the DB. Upon success, the bytes value for each entry should be set.
	ReadBatch(batch []Entry) error
}

// Write is a convenience method to perform a single entry write.
func Write(s Store, key, bytes []byte, bucket Bucket) error {
	return s.WriteBatch([]Entry{{bucket: bucket, key: key, bytes: bytes}})
}

// Read is a convenience method to perform a single entry read.
func Read(s Store, key []byte, bucket Bucket) ([]byte, error) {
	entry := Entry{bucket: bucket, key: key, bytes: nil}
	if err := s.ReadBatch([]Entry{entry}); err != nil {
		return nil, err
	}
	return entry.bytes, nil
}
