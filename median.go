package stream

import (
	heapops "container/heap"
	"errors"
)

func (s *Stats) pushMedian(x float64) {
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

// Median returns the running median of values seen.
func (s *Stats) Median() (float64, error) {
	if !s.median {
		return 0, errors.New("stream: median is not a tracked stat")
	} else if s.count == 0 {
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
