package stream

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPushMedian(t *testing.T) {
	medianStats := newMedianStats()
	for i := 0; i < 5; i++ {
		medianStats.pushMedian(float64(i))
	}

	assert.Equal(t, medianStats.lowHeap.vals, []float64{1, 0})
	assert.Equal(t, medianStats.highHeap.vals, []float64{2, 3, 4})

	medianStats.pushMedian(float64(3))

	assert.Equal(t, medianStats.lowHeap.vals, []float64{2, 0, 1})
	assert.Equal(t, medianStats.highHeap.vals, []float64{3, 3, 4})

	medianStats.pushMedian(float64(2))
	medianStats.pushMedian(float64(1))

	assert.Equal(t, medianStats.lowHeap.vals, []float64{2, 1, 1, 0})
	assert.Equal(t, medianStats.highHeap.vals, []float64{2, 3, 4, 3})
}

func TestMedian(t *testing.T) {
	medianStats := newMedianStats()
	for i := 0; i < 5; i++ {
		medianStats.pushMedian(float64(i))
	}

	median, _ := medianStats.median()

	assert.Equal(t, median, float64(2))

	medianStats.pushMedian(5)
	median, _ = medianStats.median()

	assert.Equal(t, median, float64(2.5))
}
