package stream

import (
	heapops "container/heap"
	"errors"
)

// MedianStats is a struct that keeps track of the running median.
type medianStats struct {
	lowHeap  *heap
	highHeap *heap
}

func fmax(x interface{}, y interface{}) bool {
	return x.(float64) > y.(float64)
}

func fmin(x interface{}, y interface{}) bool {
	return x.(float64) < y.(float64)
}

func newMedianStats() *medianStats {
	return &medianStats{
		lowHeap:  newHeap([]interface{}{}, fmax),
		highHeap: newHeap([]interface{}{}, fmin),
	}
}

func (s *medianStats) push(x float64) {
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

func (s *medianStats) median() (float64, error) {
	if s.lowHeap.Len()+s.highHeap.Len() == 0 {
		return 0, errors.New("stream: no values seen yet")
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
