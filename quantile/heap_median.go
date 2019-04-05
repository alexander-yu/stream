package quantile

import (
	heapops "container/heap"
	"fmt"
	"sync"

	"github.com/pkg/errors"
	"github.com/workiva/go-datastructures/queue"

	"github.com/alexander-yu/stream/quantile/heap"
)

// HeapMedian keeps track of the median of an entire stream using heaps.
type HeapMedian struct {
	window   int
	lowHeap  *heap.Heap
	highHeap *heap.Heap
	queue    *queue.RingBuffer
	mux      sync.Mutex
}

func fmax(x float64, y float64) bool {
	return x > y
}

func fmin(x float64, y float64) bool {
	return x < y
}

// NewHeapMedian instantiates a HeapMedian struct.
func NewHeapMedian(window int) (*HeapMedian, error) {
	if window < 0 {
		return nil, errors.Errorf("%d is a negative window", window)
	}

	return &HeapMedian{
		window:   window,
		lowHeap:  heap.New("low", []float64{}, fmax),
		highHeap: heap.New("high", []float64{}, fmin),
		queue:    queue.NewRingBuffer(uint64(window)),
	}, nil
}

// String returns a string representation of the metric.
func (m *HeapMedian) String() string {
	name := "quantile.HeapMedian"
	window := fmt.Sprintf("window:%v", m.window)
	return fmt.Sprintf("%s_{%s}", name, window)
}

// Push adds a number for calculating the median.
func (m *HeapMedian) Push(x float64) error {
	m.mux.Lock()
	defer m.mux.Unlock()

	var item *heap.Item
	// if queue is full, we need to remove old item
	if m.window != 0 && m.queue.Len() == uint64(m.window) {
		tail, err := m.queue.Get()
		if err != nil {
			return errors.Wrap(err, "error popping item from queue")
		}

		item = tail.(*heap.Item)
		low := item.HeapID == m.lowHeap.ID
		switch {
		// if new/old items are in the same heap, just replace the old item
		case low && x <= m.lowHeap.Peek():
			m.lowHeap.Update(item, x)
		case !low && x > m.lowHeap.Peek():
			m.highHeap.Update(item, x)
		// otherwise, remove old item and push new value to other heap
		case low && x > m.lowHeap.Peek():
			m.lowHeap.Remove(item)
			item.Val = x
			heapops.Push(m.highHeap, item)
		default:
			m.highHeap.Remove(item)
			item.Val = x
			heapops.Push(m.lowHeap, item)
		}

		// rebalance heaps
		if m.lowHeap.Len()+1 < m.highHeap.Len() {
			item = heapops.Pop(m.highHeap).(*heap.Item)
			heapops.Push(m.lowHeap, item)
		} else if m.lowHeap.Len() > m.highHeap.Len()+1 {
			item = heapops.Pop(m.lowHeap).(*heap.Item)
			heapops.Push(m.highHeap, item)
		}
	} else {
		item = &heap.Item{Val: x}
		if m.lowHeap.Len() == 0 || x <= m.lowHeap.Peek() {
			heapops.Push(m.lowHeap, item)
		} else {
			heapops.Push(m.highHeap, item)
		}

		// rebalance heaps
		if m.lowHeap.Len()+1 < m.highHeap.Len() {
			item = heapops.Pop(m.highHeap).(*heap.Item)
			heapops.Push(m.lowHeap, item)
		} else if m.lowHeap.Len() > m.highHeap.Len()+1 {
			item = heapops.Pop(m.lowHeap).(*heap.Item)
			heapops.Push(m.highHeap, item)
		}
	}

	if m.window != 0 {
		err := m.queue.Put(item)
		if err != nil {
			return errors.Wrapf(err, "error pushing %f to queue", x)
		}
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

	// return top of the larger of the heaps
	if m.lowHeap.Len() < m.highHeap.Len() {
		return m.highHeap.Peek(), nil
	} else if m.lowHeap.Len() > m.highHeap.Len() {
		return m.lowHeap.Peek(), nil
	} else {
		// otherwise, return average of the two tops
		low := m.lowHeap.Peek()
		high := m.highHeap.Peek()
		return (low + high) / 2, nil
	}
}

// Clear resets the metric.
func (m *HeapMedian) Clear() {
	m.mux.Lock()
	defer m.mux.Unlock()
	m.queue.Dispose()
	m.queue = queue.NewRingBuffer(uint64(m.window))
	m.lowHeap = heap.New("low", []float64{}, fmax)
	m.highHeap = heap.New("high", []float64{}, fmin)
}
