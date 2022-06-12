package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCompressiveString(t *testing.T) {
	rnd := NewRandom(1)
	{
		sz := 10
		result := CompressibleString(rnd, sz, 0.3)
		assert.True(t, len(result) == sz)
	}
	{
		sz := 10
		result := CompressibleString(rnd, sz, 0)
		assert.True(t, len(result) == sz)
	}

}
