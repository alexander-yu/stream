package minmax

import (
	"fmt"
	"math"
	"sync"

	"github.com/gammazero/deque"
	"github.com/pkg/errors"
	"github.com/workiva/go-datastructures/queue"
)

// Max keeps track of the maximum of a stream.
type Max struct {
	window int
	mux    sync.Mutex
	// Used if window > 0
	queue *queue.RingBuffer
	deque *deque.Deque[float64]
	// Used if window == 0
	max   float64
	count int
}

// NewMax instantiates a Max struct.
func NewMax(window int) (*Max, error) {
	if window < 0 {
		return nil, errors.Errorf("%d is a negative window", window)
	}

	return &Max{
		queue:  queue.NewRingBuffer(uint64(window)),
		deque:  deque.New[float64](),
		max:    math.Inf(-1),
		window: window,
	}, nil
}

// NewGlobalMax instantiates a global Max struct.
// This is equivalent to calling NewMax(0).
func NewGlobalMax() *Max {
	return &Max{
		queue:  queue.NewRingBuffer(uint64(0)),
		deque:  deque.New[float64](),
		max:    math.Inf(-1),
		window: 0,
	}
}

// String returns a string representation of the metric.
func (m *Max) String() string {
	name := "minmax.Max"
	window := fmt.Sprintf("window:%v", m.window)
	return fmt.Sprintf("%s_{%s}", name, window)
}

// Push adds a number for calculating the maximum.
func (m *Max) Push(x float64) error {
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

		for m.deque.Len() > 0 && m.deque.Back() < x {
			m.deque.PopBack()
		}
		m.deque.PushBack(x)

	} else {
		m.count++
		m.max = math.Max(m.max, x)
	}

	return nil
}

// Value returns the value of the maximum.
func (m *Max) Value() (float64, error) {
	m.mux.Lock()
	defer m.mux.Unlock()

	if m.count == 0 {
		return 0, errors.New("no values seen yet")
	} else if m.window == 0 {
		return m.max, nil
	}

	return m.deque.Front(), nil
}

// Clear resets the metric.
func (m *Max) Clear() {
	m.mux.Lock()
	defer m.mux.Unlock()
	m.count = 0
	m.max = math.Inf(-1)
	m.queue.Dispose()
	m.queue = queue.NewRingBuffer(uint64(m.window))
	m.deque = deque.New[float64]()
}
