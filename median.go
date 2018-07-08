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

func newMedianStats() *medianStats {
	return &medianStats{
		lowHeap:  newHeap([]float64{}, true),
		highHeap: newHeap([]float64{}, false),
	}
}

func (s *medianStats) pushMedian(x float64) {
	if s.lowHeap.Len() == 0 || x <= s.lowHeap.peek() {
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
		return s.highHeap.peek(), nil
	} else if s.lowHeap.Len() > s.highHeap.Len() {
		return s.lowHeap.peek(), nil
	} else {
		low := s.lowHeap.peek()
		high := s.highHeap.peek()
		return (low + high) / 2, nil
	}
}
