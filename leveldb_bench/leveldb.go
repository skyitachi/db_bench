package leveldb_bench

import (
	"db_bench/utils"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

type DB struct {
	db_ *leveldb.DB
}

func (db *DB) Put(key, value []byte) error {
	return db.db_.Put(key, value, nil)
}

func NewLevelDB() *DB {
	dbPath := "/tmp/leveldb_bench"
	utils.CleanDB(dbPath)
	opts := &opt.Options{}
	db, err := leveldb.OpenFile(dbPath, opts)
	if err != nil {
		panic(err)
	}

	return &DB{
		db_: db,
	}
}
