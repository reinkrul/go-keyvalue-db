package spi

type DataStore interface {
	Set(bucket string, key string, value string) error
	Get(bucket string, key string) (string, error)
	Close() error
}