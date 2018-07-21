package stream

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMedianPush(t *testing.T) {
	medianStats := newMedianStats()
	for i := 0.; i < 5; i++ {
		medianStats.push(i)
	}

	assert.Equal(t, medianStats.lowHeap.vals, []interface{}{1., 0.})
	assert.Equal(t, medianStats.highHeap.vals, []interface{}{2., 3., 4.})

	medianStats.push(3.)

	assert.Equal(t, medianStats.lowHeap.vals, []interface{}{2., 0., 1.})
	assert.Equal(t, medianStats.highHeap.vals, []interface{}{3., 3., 4.})

	medianStats.push(2.)
	medianStats.push(1.)

	assert.Equal(t, medianStats.lowHeap.vals, []interface{}{2., 1., 1., 0.})
	assert.Equal(t, medianStats.highHeap.vals, []interface{}{2., 3., 4., 3.})
}

func TestMedian(t *testing.T) {
	medianStats := newMedianStats()
	for i := 0.; i < 5; i++ {
		medianStats.push(i)
	}

	median, _ := medianStats.median()

	assert.Equal(t, median, 2.)

	medianStats.push(5)
	median, _ = medianStats.median()

	assert.Equal(t, median, 2.5)
}
