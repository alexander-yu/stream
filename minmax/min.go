package minmax

import (
	"fmt"
	"math"
	"sync"

	"github.com/gammazero/deque"
	"github.com/pkg/errors"
	"github.com/workiva/go-datastructures/queue"
)

// Min keeps track of the minimum of a stream.
type Min struct {
	window int
	mux    sync.Mutex
	count  int
	// Used if window > 0
	queue *queue.RingBuffer
	deque *deque.Deque[float64]
	// Used if window == 0
	min float64
}

// NewMin instantiates a Min struct.
func NewMin(window int) (*Min, error) {
	if window < 0 {
		return nil, errors.Errorf("%d is a negative window", window)
	}

	return &Min{
		queue:  queue.NewRingBuffer(uint64(window)),
		deque:  new(deque.Deque[float64]),
		min:    math.Inf(1),
		window: window,
	}, nil
}

// NewGlobalMin instantiates a global Min struct.
// This is equivalent to calling NewMin(0).
func NewGlobalMin() *Min {
	return &Min{
		queue:  queue.NewRingBuffer(uint64(0)),
		deque:  new(deque.Deque[float64]),
		min:    math.Inf(1),
		window: 0,
	}
}

// String returns a string representation of the metric.
func (m *Min) String() string {
	name := "minmax.Min"
	window := fmt.Sprintf("window:%v", m.window)
	return fmt.Sprintf("%s_{%s}", name, window)
}

// Push adds a number for calculating the minimum.
func (m *Min) Push(x float64) error {
	m.mux.Lock()
	defer m.mux.Unlock()

	if m.window != 0 {
		if m.queue.Len() == uint64(m.window) {
			val, err := m.queue.Get()
			if err != nil {
				return errors.Wrap(err, "error popping item from queue")
			}

			m.count--

			if m.deque.Front() == *val.(*float64) {
				m.deque.PopFront()
			}
		}

		err := m.queue.Put(&x)
		if err != nil {
			return errors.Wrapf(err, "error pushing %f to queue", x)
		}

		m.count++

		for m.deque.Len() > 0 && m.deque.Back() > x {
			m.deque.PopBack()
		}
		m.deque.PushBack(x)

	} else {
		m.count++
		m.min = math.Min(m.min, x)
	}

	return nil
}

// Value returns the value of the minimum.
func (m *Min) Value() (float64, error) {
	m.mux.Lock()
	defer m.mux.Unlock()

	if m.count == 0 {
		return 0, errors.New("no values seen yet")
	} else if m.window == 0 {
		return m.min, nil
	}

	return m.deque.Front(), nil
}

// Clear resets the metric.
func (m *Min) Clear() {
	m.mux.Lock()
	defer m.mux.Unlock()
	m.count = 0
	m.min = math.Inf(1)
	m.queue.Dispose()
	m.queue = queue.NewRingBuffer(uint64(m.window))
	m.deque = new(deque.Deque[float64])
}
