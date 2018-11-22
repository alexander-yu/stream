package median

import (
	heapops "container/heap"
	"errors"
)

// HeapMedian keeps track of the running median of an entire stream using heaps.
type HeapMedian struct {
	lowHeap  *heap
	highHeap *heap
}

func fmax(x interface{}, y interface{}) bool {
	return x.(float64) > y.(float64)
}

func fmin(x interface{}, y interface{}) bool {
	return x.(float64) < y.(float64)
}

// NewHeapMedian returns an initialized HeapMedian struct.
func NewHeapMedian() *HeapMedian {
	return &HeapMedian{
		lowHeap:  newHeap([]interface{}{}, fmax),
		highHeap: newHeap([]interface{}{}, fmin),
	}
}

// Push consumes a number from a stream for calculating the running median.
func (s *HeapMedian) Push(x float64) {
	if s.lowHeap.Len() == 0 || x <= s.lowHeap.peek().(float64) {
		heapops.Push(s.lowHeap, x)
	} else {
		heapops.Push(s.highHeap, x)
	}

	if s.lowHeap.Len()+1 < s.highHeap.Len() {
		heapops.Push(s.lowHeap, heapops.Pop(s.highHeap))
	} else if s.lowHeap.Len() > s.highHeap.Len()+1 {
		heapops.Push(s.highHeap, heapops.Pop(s.lowHeap))
	}
}

// Median returns the current running median, or error if no median is available.
func (s *HeapMedian) Median() (float64, error) {
	if s.lowHeap.Len()+s.highHeap.Len() == 0 {
		return 0, errors.New("no values seen yet")
	}

	if s.lowHeap.Len() < s.highHeap.Len() {
		return s.highHeap.peek().(float64), nil
	} else if s.lowHeap.Len() > s.highHeap.Len() {
		return s.lowHeap.peek().(float64), nil
	} else {
		low := s.lowHeap.peek().(float64)
		high := s.highHeap.peek().(float64)
		return (low + high) / 2, nil
	}
}
