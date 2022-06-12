package utils

import "fmt"

type RandomGenerator struct {
	data []byte
	pos  int
}

func (rg *RandomGenerator) Generate(sz int) []byte {
	if rg.pos+sz > len(rg.data) {
		rg.pos = 0
		if sz > len(rg.data) {
			panic(fmt.Errorf("request size %d bigger than generator size %d", sz, len(rg.data)))
		}
	}
	ret := rg.data[rg.pos : rg.pos+sz]
	rg.pos += sz
	return ret
}

func NewRandomGenerator() *RandomGenerator {
	rnd := NewRandom(301)
	cnt := 0
	data := make([]byte, 10488576)
	for cnt < 10488576 {
		ret := CompressibleString(rnd, 100, 0.8)
		copy(data[cnt:], ret)
		cnt += 100
	}
	pos := 0
	return &RandomGenerator{
		data: data,
		pos:  pos,
	}
}
