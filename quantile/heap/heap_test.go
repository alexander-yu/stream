package heap

import (
	heapops "container/heap"
	"testing"

	"github.com/stretchr/testify/assert"
)

func fmax(x float64, y float64) bool {
	return x > y
}

func fmin(x float64, y float64) bool {
	return x < y
}

func TestHeap(t *testing.T) {
	heap := NewHeap([]float64{1, 2, 3, 4}, fmax)
	assert.Equal(t, 4., heapops.Pop(heap).(*Item).Val)

	heap = NewHeap([]float64{1, 2, 3, 4}, fmin)
	assert.Equal(t, 1., heapops.Pop(heap).(*Item).Val)
}

func TestPeek(t *testing.T) {
	heap := NewHeap([]float64{1, 2, 3, 4}, fmax)
	heapops.Push(heap, &Item{Val: 5})
	heapops.Push(heap, &Item{Val: 4})
	assert.Equal(t, 5., heap.Peek())
}

func TestValues(t *testing.T) {
	heap := NewHeap([]float64{1, 2, 3, 4}, fmax)
	heapops.Push(heap, &Item{Val: 5})
	heapops.Push(heap, &Item{Val: 4})
	assert.Equal(t, []float64{5, 4, 4, 1, 2, 3}, heap.Values())
}

func TestUpdate(t *testing.T) {
	heap := NewHeap([]float64{1, 2, 3, 4}, fmax)
	item := &Item{Val: 5}

	heapops.Push(heap, item)
	assert.Equal(t, []float64{5, 4, 3, 1, 2}, heap.Values())

	heap.Update(item, 2)
	assert.Equal(t, []float64{4, 2, 3, 1, 2}, heap.Values())
}
