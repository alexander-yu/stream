package minmax

import (
	"math"
	"sync"

	"github.com/gammazero/deque"
	"github.com/pkg/errors"
	"github.com/workiva/go-datastructures/queue"

	"github.com/alexander-yu/stream"
)

// Min keeps track of the running minimum of a stream. Note: for global minimums,
// Core.Min() also tracks this; Min provides the additional functionality of keeping
// track of minimums over a rolling window.
type Min struct {
	window int
	core   *stream.Core
	mux    sync.Mutex
	// Used if window > 0
	queue *queue.RingBuffer
	deque *deque.Deque
	// Used if window == 0
	min   float64
	count int
}

// NewMin instantiates a Min struct.
func NewMin(window int) (*Min, error) {
	if window < 0 {
		return nil, errors.Errorf("%d is a negative window", window)
	}

	return &Min{
		queue:  queue.NewRingBuffer(uint64(window)),
		deque:  &deque.Deque{},
		min:    math.Inf(1),
		window: window,
	}, nil
}

// Subscribe subscribes the Min to a Core object.
func (m *Min) Subscribe(c *stream.Core) {
	m.core = c
}

// Config returns the CoreConfig needed.
func (m *Min) Config() *stream.CoreConfig {
	return &stream.CoreConfig{
		Window: &m.window,
	}
}

// Push adds a number for calculating the running minimum.
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

			if m.deque.Front().(*float64) == val.(*float64) {
				m.deque.PopFront()
			}
		}

		err := m.queue.Put(&x)
		if err != nil {
			return errors.Wrapf(err, "error pushing %f to queue", x)
		}

		m.count++

		for m.deque.Len() > 0 && *m.deque.Back().(*float64) > x {
			m.deque.PopBack()
		}
		m.deque.PushBack(&x)

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

	return *m.deque.Front().(*float64), nil
}
