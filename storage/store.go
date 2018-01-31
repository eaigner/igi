package storage

const (
	TransactionBucket Bucket = 1
)

var allBuckets = []Bucket{
	TransactionBucket,
}

var bucketKeys = map[Bucket][]byte{}

func init() {
	// Converty bucket int ids to []byte keys
	for _, bucket := range allBuckets {
		bucketKeys[bucket] = []byte{byte(bucket)}
	}
}

type Bucket byte

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

	// Close closes the store
	Close() error
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

// Exists returns true if the key exists in the provided bucket
func Exists(s Store, key []byte, bucket Bucket) (bool, error) {
	item, err := Read(s, key, bucket)
	if err != nil {
		return false, err
	}
	return item != nil, nil
}
