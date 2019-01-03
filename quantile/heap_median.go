package quantile

import (
	heapops "container/heap"
	"errors"
	"sync"

	"github.com/alexander-yu/stream/quantile/heap"
)

// HeapMedian keeps track of the median of an entire stream using heaps.
type HeapMedian struct {
	lowHeap  *heap.Heap
	highHeap *heap.Heap
	mux      sync.Mutex
}

func fmax(x float64, y float64) bool {
	return x > y
}

func fmin(x float64, y float64) bool {
	return x < y
}

// NewHeapMedian instantiates a HeapMedian struct.
func NewHeapMedian() *HeapMedian {
	return &HeapMedian{
		lowHeap:  heap.NewHeap([]float64{}, fmax),
		highHeap: heap.NewHeap([]float64{}, fmin),
	}
}

// Push adds a number for calculating the median.
func (m *HeapMedian) Push(x float64) error {
	m.mux.Lock()
	defer m.mux.Unlock()

	item := &heap.Item{Val: x}
	if m.lowHeap.Len() == 0 || x <= m.lowHeap.Peek() {
		heapops.Push(m.lowHeap, item)
	} else {
		heapops.Push(m.highHeap, item)
	}

	if m.lowHeap.Len()+1 < m.highHeap.Len() {
		heapops.Push(m.lowHeap, heapops.Pop(m.highHeap))
	} else if m.lowHeap.Len() > m.highHeap.Len()+1 {
		heapops.Push(m.highHeap, heapops.Pop(m.lowHeap))
	}

	return nil
}

// Value returns the value of the median.
func (m *HeapMedian) Value() (float64, error) {
	m.mux.Lock()
	defer m.mux.Unlock()

	if m.lowHeap.Len()+m.highHeap.Len() == 0 {
		return 0, errors.New("no values seen yet")
	}

	if m.lowHeap.Len() < m.highHeap.Len() {
		return m.highHeap.Peek(), nil
	} else if m.lowHeap.Len() > m.highHeap.Len() {
		return m.lowHeap.Peek(), nil
	} else {
		low := m.lowHeap.Peek()
		high := m.highHeap.Peek()
		return (low + high) / 2, nil
	}
}
