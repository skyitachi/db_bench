package boltdb_bench

import (
	"db_bench/utils"
	bolt "go.etcd.io/bbolt"
)

type DB struct {
	db     *bolt.DB
	bucket string
}

func (db *DB) Put(key, value []byte) error {
	return db.db.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte(db.bucket))
		return bkt.Put(key, value)
	})
}

func (db *DB) Begin() (*bolt.Tx, error) {
	return db.db.Begin(true)
}

func (db *DB) GetBucket(tx *bolt.Tx) *bolt.Bucket {
	return tx.Bucket([]byte(db.bucket))
}

func NewBoltDB() *DB {
	dbPath := "/tmp/boltdb_bench"
	utils.CleanDB(dbPath)

	db, err := bolt.Open(dbPath, 0666, &bolt.Options{
		NoSync: true,
	})

	if err != nil {
		panic(err)
	}

	bucketName := "boltdb_bench"

	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		return err
	}); err != nil {
		panic(err)
	}

	return &DB{
		db:     db,
		bucket: bucketName,
	}
}
