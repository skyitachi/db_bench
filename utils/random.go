package utils

const M uint32 = 2147483647
const A uint64 = 16807

type Random struct {
	seed uint32
}

func (r *Random) Next() uint32 {
	var product uint64
	product = uint64(r.seed) * A
	r.seed = uint32(product>>31 + (product & uint64(M)))

	if r.seed > M {
		r.seed -= M
	}
	return r.seed
}

func (r *Random) Uniform(n int) uint32 {
	return r.Next() % uint32(n)
}

func (r *Random) OneIn(n int) bool {
	return (r.Next() % uint32(n)) == 0
}

func (r *Random) Skewed(maxLog int) uint32 {
	return r.Uniform(1 << r.Uniform(maxLog+1))
}

func NewRandom(seed uint32) *Random {
	seed = seed & 0x7fffffff
	if seed == 0 || seed == M {
		seed = 1
	}
	return &Random{
		seed: seed,
	}
}
