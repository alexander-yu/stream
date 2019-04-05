package quantile

import (
	"fmt"
	"math"
	"strings"
	"sync"

	"github.com/alexander-yu/stream/quantile/order"
	"github.com/pkg/errors"
	"github.com/workiva/go-datastructures/queue"
)

// Quantile keeps track of the quantile of a stream using order statistics.
type Quantile struct {
	window        int
	interpolation Interpolation
	queue         *queue.RingBuffer
	statistic     order.Statistic
	mux           sync.RWMutex
}

// NewQuantile instantiates a Quantile struct.
func NewQuantile(window int, options ...Option) (*Quantile, error) {
	if window < 0 {
		return nil, errors.Errorf("attempted to set negative window of %d", window)
	}

	avl, err := AVL.init()
	if err != nil {
		return nil, errors.Wrap(err, "error instantiating default AVL order.Statistic")
	}

	quantile := &Quantile{
		window:        window,
		interpolation: Linear,
		queue:         queue.NewRingBuffer(uint64(window)),
		statistic:     avl,
	}

	for _, option := range options {
		err = option(quantile)
		if err != nil {
			return nil, errors.Wrap(err, "error setting option")
		}
	}

	return quantile, nil
}

// String returns a string representation of the metric.
func (q *Quantile) String() string {
	name := "quantile.Quantile"
	params := []string{
		fmt.Sprintf("window:%v", q.window),
		fmt.Sprintf("interpolation:%v", q.interpolation),
	}
	return fmt.Sprintf("%s_{%s}", name, strings.Join(params, ","))
}

// Push adds a number for calculating the quantile.
func (q *Quantile) Push(x float64) error {
	q.mux.Lock()
	defer q.mux.Unlock()

	if q.window != 0 {
		if q.queue.Len() == uint64(q.window) {
			val, err := q.queue.Get()
			if err != nil {
				return errors.Wrap(err, "error popping item from queue")
			}

			y := val.(float64)
			q.statistic.Remove(y)
		}

		err := q.queue.Put(x)
		if err != nil {
			return errors.Wrapf(err, "error pushing %f to queue", x)
		}
	}

	q.statistic.Add(x)
	return nil
}

// Value returns the value of the quantile.
func (q *Quantile) Value(quantile float64) (float64, error) {
	if quantile <= 0 || quantile >= 1 {
		return 0, errors.Errorf("quantile %f not in (0, 1)", quantile)
	}

	q.mux.RLock()
	defer q.mux.RUnlock()

	size := int(q.statistic.Size())
	if size == 0 {
		return 0, errors.New("no values seen yet")
	}

	idxRaw := quantile * float64(size-1)
	idxTrunc := math.Trunc(idxRaw)
	idx := int(idxTrunc)
	// if the estimated index is actually an integer,
	// no interpolation needed
	if idxRaw == idxTrunc {
		return q.statistic.Select(idx).Value(), nil
	}

	delta := idxRaw - idxTrunc
	switch q.interpolation {
	case Linear:
		lo := q.statistic.Select(idx).Value()
		hi := q.statistic.Select(idx + 1).Value()
		return (1-delta)*lo + delta*hi, nil
	case Lower:
		return q.statistic.Select(idx).Value(), nil
	case Higher:
		return q.statistic.Select(idx + 1).Value(), nil
	case Nearest:
		switch {
		case delta == 0.5:
			if idx%2 == 0 {
				return q.statistic.Select(idx).Value(), nil
			}
			return q.statistic.Select(idx + 1).Value(), nil
		case delta < 0.5:
			return q.statistic.Select(idx).Value(), nil
		default:
			return q.statistic.Select(idx + 1).Value(), nil
		}
	default:
		lo := q.statistic.Select(idx).Value()
		hi := q.statistic.Select(idx + 1).Value()
		return (lo + hi) / 2., nil
	}
}

// Clear resets the metric.
func (q *Quantile) Clear() {
	q.mux.Lock()
	defer q.mux.Unlock()
	q.queue.Dispose()
	q.queue = queue.NewRingBuffer(uint64(q.window))
	q.statistic.Clear()
}

// RLock locks the quantile for reading.
func (q *Quantile) RLock() {
	q.mux.RLock()
}

// RUnlock undoes an RLock call.
func (q *Quantile) RUnlock() {
	q.mux.RUnlock()
}
