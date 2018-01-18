package storage

type Store interface {
	Set(bytes []byte, key, domain string) bool
	Get(key, domain string) ([]byte, bool)
}
