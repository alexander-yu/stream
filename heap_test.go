package stream

import (
	heapops "container/heap"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeap(t *testing.T) {
	heap := newHeap([]float64{1, 2, 3, 4}, true)
	assert.Equal(t, float64(4), heapops.Pop(heap))

	heap = newHeap([]float64{1, 2, 3, 4}, false)
	assert.Equal(t, float64(1), heapops.Pop(heap))
}

func TestPeek(t *testing.T) {
	heap := newHeap([]float64{1, 2, 3, 4}, true)
	heapops.Push(heap, float64(5))
	heapops.Push(heap, float64(4))
	assert.Equal(t, float64(5), heap.peek())
}
