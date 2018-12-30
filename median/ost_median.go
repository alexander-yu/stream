package median

import (
	"sync"

	"github.com/pkg/errors"
	"github.com/workiva/go-datastructures/queue"

	"github.com/alexander-yu/stream/median/ost"
)

// OSTMedian keeps track of the median of a stream using order statistic trees.
type OSTMedian struct {
	queue  *queue.RingBuffer
	tree   ost.Tree
	window int
	mux    sync.Mutex
}

// NewOSTMedian instantiates an OSTMedian struct. The implementation of the
// underlying order statistic tree can be configured by passing in a constant
// of type ost.Impl.
func NewOSTMedian(window int, impl ost.Impl) (*OSTMedian, error) {
	if window < 0 {
		return nil, errors.Errorf("%d is a negative window", window)
	}

	tree, err := impl.EmptyTree()
	if err != nil {
		return nil, errors.Wrap(err, "error instantiating empty ost.Tree")
	}

	return &OSTMedian{
		queue:  queue.NewRingBuffer(uint64(window)),
		tree:   tree,
		window: window,
	}, nil
}

// Push adds a number for calculating the median.
func (m *OSTMedian) Push(x float64) error {
	m.mux.Lock()
	defer m.mux.Unlock()

	if m.window != 0 {
		if m.queue.Len() == uint64(m.window) {
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
	}

	m.tree.Add(x)
	return nil
}

// Value returns the value of the median.
func (m *OSTMedian) Value() (float64, error) {
	m.mux.Lock()
	defer m.mux.Unlock()

	size := int(m.queue.Len())
	if size == 0 {
		return 0, errors.New("no values seen yet")
	} else if size%2 == 0 {
		left := m.tree.Select(size/2 - 1).Value()
		right := m.tree.Select(size / 2).Value()
		return float64(left+right) / float64(2), nil
	}

	return m.tree.Select(size / 2).Value(), nil
}
