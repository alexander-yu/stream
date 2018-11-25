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
	window *int
	core   *stream.Core
	mux    sync.Mutex
}

// NewAVLMedian instantiates an AVLMedian struct
func NewAVLMedian(window int) (*AVLMedian, error) {
	if window < 0 {
		return nil, errors.New("window is negative")
	}

	return &AVLMedian{
		queue:  queue.NewRingBuffer(uint64(window)),
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
	return nil
}

// Value returns the value of the median.
func (m *AVLMedian) Value() (float64, error) {
	m.mux.Lock()
	defer m.mux.Unlock()
	return 0, nil
}
