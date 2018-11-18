package median

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeapMedianPush(t *testing.T) {
	medianStats := NewHeapMedian()
	for i := 0.; i < 5; i++ {
		medianStats.Push(i)
	}

	assert.Equal(t, medianStats.lowHeap.vals, []interface{}{1., 0.})
	assert.Equal(t, medianStats.highHeap.vals, []interface{}{2., 3., 4.})

	medianStats.Push(3.)

	assert.Equal(t, medianStats.lowHeap.vals, []interface{}{2., 0., 1.})
	assert.Equal(t, medianStats.highHeap.vals, []interface{}{3., 3., 4.})

	medianStats.Push(2.)
	medianStats.Push(1.)

	assert.Equal(t, medianStats.lowHeap.vals, []interface{}{2., 1., 1., 0.})
	assert.Equal(t, medianStats.highHeap.vals, []interface{}{2., 3., 4., 3.})
}

func TestHeapMedian(t *testing.T) {
	medianStats := NewHeapMedian()
	for i := 0.; i < 5; i++ {
		medianStats.Push(i)
	}

	median, _ := medianStats.Median()

	assert.Equal(t, median, 2.)

	medianStats.Push(5)
	median, _ = medianStats.Median()

	assert.Equal(t, median, 2.5)
}
