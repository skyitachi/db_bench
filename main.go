package main

import (
	"db_bench/boltdb_bench"
	"db_bench/leveldb_bench"
	"db_bench/nutsdb_bench"
	"db_bench/pebbledb_bench"
	"db_bench/utils"
	"flag"
	"fmt"
	"log"
	"reflect"
	"time"
	"unsafe"
)

var gDB string
var gNum int
var gValueSize int
var gKeySize int

type KeyBuffer struct {
	buffer_ [1024]byte
	kPrefix int
}

func (kb *KeyBuffer) Set(k int) {
	ks := fmt.Sprintf("%16d", k)
	copy(kb.buffer_[kb.kPrefix:], ks)
}

func (kb KeyBuffer) Slice() []byte {
	return kb.buffer_[kb.kPrefix : kb.kPrefix+gKeySize]
}

func (kb KeyBuffer) String() string {
	bytes := kb.buffer_[kb.kPrefix : kb.kPrefix+gKeySize]
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&bytes))
	stringHeader := reflect.StringHeader{Data: sliceHeader.Data, Len: sliceHeader.Len}
	return *(*string)(unsafe.Pointer(&stringHeader))
}

func NewKeyBuffer(prefix int) *KeyBuffer {
	kb := &KeyBuffer{
		kPrefix: prefix,
	}
	for i := 0; i < prefix; i++ {
		kb.buffer_[i] = 'a'
	}
	return kb
}

func benchBoltDB() {
	db := boltdb_bench.NewBoltDB()

	gen := utils.NewRandomGenerator()

	rnd := utils.NewRandom(1)
	nums := gNum

	keyBuffer := NewKeyBuffer(0)
	valueSize := gValueSize

	var bytes int64
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	bkt := db.GetBucket(tx)

	start := time.Now()
	for i := 0; i < nums; i++ {
		k := rnd.Uniform(nums)
		keyBuffer.Set(int(k))
		if err := bkt.Put(keyBuffer.Slice(), gen.Generate(valueSize)); err != nil {
			panic(err)
		}
		bytes += int64(valueSize + len(keyBuffer.Slice()))
	}
	if err := tx.Commit(); err != nil {
		panic(err)
	}

	el := time.Since(start).Milliseconds()
	fmt.Printf("boltdb_bench consums: %d ms, bytes: %d, rate: %.2f MB/s\n", el, bytes, utils.ComputeRates(bytes, el))
}

func benchNutsDB() {
	db := nutsdb_bench.NewNutsDB()

	gen := utils.NewRandomGenerator()

	rnd := utils.NewRandom(1)
	nums := gNum

	keyBuffer := NewKeyBuffer(0)
	valueSize := gValueSize

	start := time.Now()
	tx, err := db.Begin()

	if err != nil {
		panic(err)
	}

	var bytes int64
	for i := 0; i < nums; i++ {
		k := rnd.Uniform(nums)
		keyBuffer.Set(int(k))
		if err := tx.Put(db.BucketName(), keyBuffer.Slice(), gen.Generate(valueSize), 0); err != nil {
			panic(err)
		}
		bytes += int64(valueSize + len(keyBuffer.Slice()))
	}

	if err := tx.Commit(); err != nil {
		panic(err)
	}

	el := time.Since(start).Milliseconds()
	fmt.Printf("nutsdb_bench consums: %d ms, bytes: %d, rate: %.2f MB/s\n", el, bytes, utils.ComputeRates(bytes, el))
}

func benchLevelDB() {
	db := leveldb_bench.NewLevelDB()

	gen := utils.NewRandomGenerator()

	rnd := utils.NewRandom(1)
	nums := gNum

	keyBuffer := NewKeyBuffer(0)
	valueSize := gValueSize

	start := time.Now()

	var bytes int64
	for i := 0; i < nums; i++ {
		k := rnd.Uniform(nums)
		keyBuffer.Set(int(k))
		if err := db.Put(keyBuffer.Slice(), gen.Generate(valueSize)); err != nil {
			panic(err)
		}
		bytes += int64(valueSize + len(keyBuffer.Slice()))
	}

	el := time.Since(start).Milliseconds()
	fmt.Printf("%s consums: %d ms, bytes: %d, rate: %.2f MB/s\n", "leveldb_bench", el, bytes, utils.ComputeRates(bytes, el))
}

func benchPebbleDB() {
	db := pebbledb_bench.NewPebbleDB()

	gen := utils.NewRandomGenerator()

	rnd := utils.NewRandom(1)
	nums := gNum

	keyBuffer := NewKeyBuffer(0)
	valueSize := gValueSize

	start := time.Now()

	var bytes int64
	b := db.NewBatch()
	for i := 0; i < nums; i++ {
		k := rnd.Uniform(nums)
		keyBuffer.Set(int(k))
		if err := b.Set(keyBuffer.Slice(), gen.Generate(valueSize), nil); err != nil {
			panic(err)
		}
		bytes += int64(valueSize + len(keyBuffer.Slice()))
	}

	el := time.Since(start).Milliseconds()
	if err := b.Commit(nil); err != nil {
		panic(err)
	}
	fmt.Printf("%s_bench consums: %d ms, bytes: %d, rate: %.2f MB/s\n", gDB, el, bytes, utils.ComputeRates(bytes, el))

}

func main() {
	flag.StringVar(&gDB, "db", "boltdb", "specify kv_engine name")
	flag.IntVar(&gNum, "num", 100000, "specify how many w/r calls to db")
	flag.IntVar(&gValueSize, "value_size", 100, "usage specify value_size")
	flag.IntVar(&gKeySize, "key_size", 16, "usage specify value_size")
	flag.Parse()

	switch gDB {
	case "boltdb":
		benchBoltDB()
	case "nutsdb":
		benchNutsDB()
	case "leveldb":
		benchLevelDB()
	case "pebbledb":
		benchPebbleDB()
	default:
		log.Fatalf("%s db not supported", gDB)
	}

}
