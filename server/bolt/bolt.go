package bolt

import (
	"github.com/reinkrul/go-keyvalue-db/server/spi"
	bolt "go.etcd.io/bbolt"
	"log"
)

type impl struct {
	db *bolt.DB
}

func Connect(path string) (spi.DataStore, error) {
	db, err := bolt.Open(path, 0666, nil)
	if err != nil {
		return nil, err
	}
	return impl{db: db}, nil
}

func (s impl) Set(bucket string, key string, value string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b, err := s.getBucket(tx, bucket, true)
		if err != nil {
			return err
		}
		return b.Put([]byte(key), []byte(value))
	})
}

func (s impl) Get(bucket string, key string) (string, error) {
	value := ""
	err := s.db.View(func(tx *bolt.Tx) error {
		b, err := s.getBucket(tx, bucket,false)
		if err != nil {
			return err
		}
		if b == nil {
			return nil
		}
		value = string(b.Get([]byte(key)))
		return err
	})
	return value, err
}

func (s impl) Close() error {
	log.Print("Closing Bolt database")
	return s.db.Close()
}

func (s impl) getBucket(tx *bolt.Tx, bucket string, create bool) (*bolt.Bucket, error) {
	result := tx.Bucket([]byte(bucket))
	if create && result == nil {
		return tx.CreateBucket([]byte(bucket))
	}
	return result, nil
}
