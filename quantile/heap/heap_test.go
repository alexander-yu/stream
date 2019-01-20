package heap

import (
	heapops "container/heap"
	"testing"

	"github.com/stretchr/testify/suite"
)

func fmax(x float64, y float64) bool {
	return x > y
}

func fmin(x float64, y float64) bool {
	return x < y
}

type HeapSuite struct {
	suite.Suite
	maxHeap *Heap
	minHeap *Heap
}

func TestHeapSuite(t *testing.T) {
	suite.Run(t, &HeapSuite{})
}

func (s *HeapSuite) SetupTest() {
	s.maxHeap = NewHeap("a", []float64{1, 2, 3, 4}, fmax)
	s.minHeap = NewHeap("", []float64{1, 2, 3, 4}, fmin)
}

func (s *HeapSuite) TestNewHeap() {
	s.Equal("a", s.maxHeap.ID)
	s.Equal(4., heapops.Pop(s.maxHeap).(*Item).Val)

	s.Equal("", s.minHeap.ID)
	s.Equal(1., heapops.Pop(s.minHeap).(*Item).Val)
}

func (s *HeapSuite) TestPeek() {
	heapops.Push(s.maxHeap, &Item{Val: 5})
	heapops.Push(s.maxHeap, &Item{Val: 4})
	s.Equal(5., s.maxHeap.Peek())
}

func (s *HeapSuite) TestValues() {
	heapops.Push(s.maxHeap, &Item{Val: 5})
	heapops.Push(s.maxHeap, &Item{Val: 4})
	s.Equal([]float64{5, 4, 4, 1, 2, 3}, s.maxHeap.Values())
}

func (s *HeapSuite) TestUpdate() {
	item := &Item{
		Val:    5,
		HeapID: "b",
	}

	heapops.Push(s.maxHeap, item)
	s.Equal([]float64{5, 4, 3, 1, 2}, s.maxHeap.Values())
	s.Equal("a", item.HeapID)

	s.maxHeap.Update(item, 2)
	s.Equal([]float64{4, 2, 3, 1, 2}, s.maxHeap.Values())
}

func (s *HeapSuite) TestRemove() {
	item := &Item{Val: 5}

	heapops.Push(s.maxHeap, item)
	s.Equal([]float64{5, 4, 3, 1, 2}, s.maxHeap.Values())

	s.maxHeap.Remove(item)
	s.Equal([]float64{4, 2, 3, 1}, s.maxHeap.Values())
}
