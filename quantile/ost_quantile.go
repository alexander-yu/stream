package quantile

import (
	"fmt"
	"math"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"github.com/workiva/go-datastructures/queue"

	"github.com/alexander-yu/stream/quantile/ost"
)

// OSTQuantile keeps track of the quantile of a stream using order statistic trees.
type OSTQuantile struct {
	window        int
	interpolation Interpolation
	queue         *queue.RingBuffer
	tree          ost.Tree
	mux           sync.Mutex
}

// NewOSTQuantile instantiates an OSTQuantile struct. The implementation of the
// underlying order statistic tree can be configured by passing in a constant
// of type ost.Impl.
func NewOSTQuantile(config *Config, impl ost.Impl) (*OSTQuantile, error) {
	// set defaults for any remaining unset fields
	config = setConfigDefaults(config)

	// validate config
	err := validateConfig(config)
	if err != nil {
		return nil, errors.Wrap(err, "error validating config")
	}

	tree, err := impl.EmptyTree()
	if err != nil {
		return nil, errors.Wrap(err, "error instantiating empty ost.Tree")
	}

	return &OSTQuantile{
		window:        *config.Window,
		interpolation: *config.Interpolation,
		queue:         queue.NewRingBuffer(uint64(*config.Window)),
		tree:          tree,
	}, nil
}

// String returns a string representation of the metric.
func (q *OSTQuantile) String() string {
	name := "quantile.OSTQuantile"
	params := []string{
		fmt.Sprintf("window:%v", q.window),
		fmt.Sprintf("interpolation:%v", q.interpolation),
	}
	return fmt.Sprintf("%s_{%s}", name, strings.Join(params, ","))
}

// Push adds a number for calculating the quantile.
func (q *OSTQuantile) Push(x float64) error {
	q.mux.Lock()
	defer q.mux.Unlock()

	if q.window != 0 {
		if q.queue.Len() == uint64(q.window) {
			val, err := q.queue.Get()
			if err != nil {
				return errors.Wrap(err, "error popping item from queue")
			}

			y := val.(float64)
			q.tree.Remove(y)
		}

		err := q.queue.Put(x)
		if err != nil {
			return errors.Wrapf(err, "error pushing %f to queue", x)
		}
	}

	q.tree.Add(x)
	return nil
}

// Value returns the value of the quantile.
func (q *OSTQuantile) Value(quantile float64) (float64, error) {
	if quantile <= 0 || quantile >= 1 {
		return 0, errors.Errorf("quantile %f not in (0, 1)", quantile)
	}

	q.mux.Lock()
	defer q.mux.Unlock()

	size := int(q.tree.Size())
	if size == 0 {
		return 0, errors.New("no values seen yet")
	}

	idxRaw := quantile * float64(size-1)
	idxTrunc := math.Trunc(idxRaw)
	idx := int(idxTrunc)
	// if the estimated index is actually an integer,
	// no interpolation needed
	if idxRaw == idxTrunc {
		return q.tree.Select(idx).Value(), nil
	}

	delta := idxRaw - idxTrunc
	switch q.interpolation {
	case Linear:
		lo := q.tree.Select(idx).Value()
		hi := q.tree.Select(idx + 1).Value()
		return (1-delta)*lo + delta*hi, nil
	case Lower:
		return q.tree.Select(idx).Value(), nil
	case Higher:
		return q.tree.Select(idx + 1).Value(), nil
	case Nearest:
		switch {
		case delta == 0.5:
			if idx%2 == 0 {
				return q.tree.Select(idx).Value(), nil
			}
			return q.tree.Select(idx + 1).Value(), nil
		case delta < 0.5:
			return q.tree.Select(idx).Value(), nil
		default:
			return q.tree.Select(idx + 1).Value(), nil
		}
	default:
		lo := q.tree.Select(idx).Value()
		hi := q.tree.Select(idx + 1).Value()
		return (lo + hi) / 2., nil
	}
}

// Clear resets the metric.
func (q *OSTQuantile) Clear() {
	q.mux.Lock()
	defer q.mux.Unlock()
	q.queue.Dispose()
	q.queue = queue.NewRingBuffer(uint64(q.window))
	q.tree.Clear()
}
