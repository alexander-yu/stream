package stream

import (
	heapops "container/heap"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeap(t *testing.T) {
	heap := newHeap([]float64{1, 2, 3, 4}, true)
	assert.Equal(t, 4., heapops.Pop(heap))

	heap = newHeap([]float64{1, 2, 3, 4}, false)
	assert.Equal(t, 1., heapops.Pop(heap))
}

func TestPeek(t *testing.T) {
	heap := newHeap([]float64{1, 2, 3, 4}, true)
	heapops.Push(heap, 5.)
	heapops.Push(heap, 4.)
	assert.Equal(t, 5., heap.peek())
}
