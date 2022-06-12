package nutsdb_bench

import (
	"db_bench/utils"
	"github.com/xujiajun/nutsdb"
)

type DB struct {
	db     *nutsdb.DB
	bucket string
	ttl    uint32
}

func (db *DB) Put(key, value []byte) error {
	return db.db.Update(func(tx *nutsdb.Tx) error {
		return tx.Put(db.bucket, key, value, db.ttl)
	})
}

func (db *DB) Begin() (*nutsdb.Tx, error) {
	return db.db.Begin(true)
}

func (db DB) BucketName() string {
	return db.bucket
}

func NewNutsDB() *DB {
	opt := nutsdb.DefaultOptions
	opt.Dir = "/tmp/nutsdb"

	utils.CleanDB(opt.Dir)

	db, err := nutsdb.Open(opt)
	if err != nil {
		panic(err)
	}
	return &DB{
		db:     db,
		bucket: "test_bucket",
		ttl:    0,
	}
}
