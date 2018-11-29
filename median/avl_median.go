package median

import (
	"sync"

	"github.com/pkg/errors"
	"github.com/workiva/go-datastructures/queue"

	"github.com/alexander-yu/stream"
)

// AVLMedian keeps track of the running median of a stream using AVL trees.
type AVLMedian struct {
	queue  *queue.RingBuffer
	tree   *OrderStatisticTree
	window *int
	core   *stream.Core
	mux    sync.Mutex
}

// NewAVLMedian instantiates an AVLMedian struct.
func NewAVLMedian(window int) (*AVLMedian, error) {
	if window < 0 {
		return nil, errors.Errorf("%d is a negative window", window)
	}

	return &AVLMedian{
		queue:  queue.NewRingBuffer(uint64(window)),
		tree:   &OrderStatisticTree{},
		window: stream.IntPtr(window),
	}, nil
}

// Subscribe subscribes the AVLMedian to a Core object.
func (m *AVLMedian) Subscribe(c *stream.Core) {
	m.core = c
}

// Config returns the CoreConfig needed.
func (m *AVLMedian) Config() *stream.CoreConfig {
	return &stream.CoreConfig{
		Window:      m.window,
		PushMetrics: []stream.Metric{m},
	}
}

// Push adds a number for calculating the running median.
func (m *AVLMedian) Push(x float64) error {
	m.mux.Lock()
	defer m.mux.Unlock()

	if m.queue.Len() == uint64(*m.window) {
		val, err := m.queue.Get()
		if err != nil {
			return errors.Wrap(err, "error popping item from queue")
		}

		y := val.(float64)
		m.tree.Remove(y)
	}

	err := m.queue.Put(x)
	if err != nil {
		return errors.Wrapf(err, "error pushing %f to queue", x)
	}
	m.tree.Add(x)

	return nil
}

// Value returns the value of the median.
func (m *AVLMedian) Value() (float64, error) {
	m.mux.Lock()
	defer m.mux.Unlock()

	size := int(m.queue.Len())
	if size%2 == 0 {
		left := m.tree.Select(size/2 - 1).val
		right := m.tree.Select(size / 2).val
		return float64(left+right) / float64(2), nil
	}

	return m.tree.Select(size / 2).val, nil
}
