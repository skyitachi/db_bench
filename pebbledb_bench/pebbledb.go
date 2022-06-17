package pebbledb_bench

import (
	"db_bench/utils"
	"github.com/cockroachdb/pebble"
)

type DB struct {
	db_ *pebble.DB
}

func (db *DB) NewBatch() *pebble.Batch {
	return db.db_.NewBatch()
}

func (db *DB) Put(key, value []byte) error {
	return db.db_.Set(key, value, nil)
}

func NewPebbleDB() *DB {
	dbPath := "/tmp/pebbledb_bench"
	utils.CleanDB(dbPath)
	db, err := pebble.Open(dbPath, &pebble.Options{})
	if err != nil {
		panic(err)
	}

	return &DB{
		db_: db,
	}
}
