package sample

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testReservoir() *Reservoir {
	r := &Reservoir{Size: 3}
	rand.Seed(1)

	for i := 0; i < 10; i++ {
		r.Push(i)
	}

	return r
}

func TestReservoirPush(t *testing.T) {
	r := testReservoir()
	assert.Equal(t, []interface{}{6, 7, 4}, r.sample)
}

func TestSample(t *testing.T) {
	r := testReservoir()
	sample := r.Sample()

	assert.Equal(t, []interface{}{6, 7, 4}, sample)

	sample[0] = 1

	assert.Equal(t, []interface{}{6, 7, 4}, r.sample)
}
