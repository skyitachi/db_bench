package utils

import (
	"github.com/xujiajun/utils/filesystem"
	"os"
)

var kTestChars = []byte{'\x00', '\x01', 'a', 'b', 'c', 'd', 'e', '\xfd', '\xfe', '\xff'}

func RandomBytes(rnd *Random, sz int) []byte {
	ret := make([]byte, sz)

	for i := 0; i < sz; i++ {
		ret[i] = ' ' + uint8(rnd.Uniform(95))
	}
	return ret
}

func RandomKey(rnd *Random, sz int) []byte {
	ret := make([]byte, sz)

	for i := 0; i < sz; i++ {
		ret[i] = kTestChars[rnd.Uniform(len(kTestChars))]
	}
	return ret
}

func CompressibleString(rnd *Random, sz int, compressed_fraction float64) []byte {
	raw := int(float64(sz) * compressed_fraction)
	if raw < 1 {
		raw = 1
	}
	rawData := RandomBytes(rnd, raw)
	result := make([]byte, sz)

	cnt := 0
	for cnt < sz {
		copy(result[cnt:], rawData)
		cnt += raw
	}
	return result
}

func ComputeRates(bytes, mills int64) float64 {
	bf64 := float64(bytes)
	mf64 := float64(mills)

	return (bf64 / 1024.0 / 1024.0) / (mf64 / 1000.0)
}

func CleanDB(path string) {

	if !filesystem.PathIsExist(path) {
		return
	}
	if err := os.RemoveAll(path); err != nil {
		panic(err)
	}
}
