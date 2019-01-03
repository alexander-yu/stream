package heap

import (
	heapops "container/heap"
	"testing"

	"github.com/stretchr/testify/assert"
)

func imax(x interface{}, y interface{}) bool {
	return x.(int) > y.(int)
}

func imin(x interface{}, y interface{}) bool {
	return x.(int) < y.(int)
}

func TestHeap(t *testing.T) {
	heap := NewHeap([]interface{}{1, 2, 3, 4}, imax)
	assert.Equal(t, 4, heapops.Pop(heap))

	heap = NewHeap([]interface{}{1, 2, 3, 4}, imin)
	assert.Equal(t, 1, heapops.Pop(heap))
}

func TestPeek(t *testing.T) {
	heap := NewHeap([]interface{}{1, 2, 3, 4}, imax)
	heapops.Push(heap, 5)
	heapops.Push(heap, 4)
	assert.Equal(t, 5, heap.Peek())
}

func TestValues(t *testing.T) {
	heap := NewHeap([]interface{}{1, 2, 3, 4}, imax)
	heapops.Push(heap, 5)
	heapops.Push(heap, 4)
	assert.Equal(t, []interface{}{5, 4, 4, 1, 2, 3}, heap.Values())
}
