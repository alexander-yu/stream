package joint

import (
	"github.com/pkg/errors"
	"github.com/workiva/go-datastructures/queue"
)

// Autocorrelation is a metric that tracks the sample autocorrelation.
type Autocorrelation struct {
	lag         uint64
	queue       *queue.RingBuffer
	correlation *Correlation
	core        *Core
}

// NewAutocorrelation instantiates an Autocorrelation struct.
func NewAutocorrelation(lag int, window int) (*Autocorrelation, error) {
	if lag < 1 {
		return nil, errors.Errorf("%d is a lag that is less than 1", lag)
	}

	correlation, err := NewCorrelation(window)
	if err != nil {
		return nil, errors.Wrap(err, "error creating Correlation")
	}

	autocorrelation := &Autocorrelation{
		lag:         uint64(lag),
		queue:       queue.NewRingBuffer(uint64(lag)),
		correlation: correlation,
	}

	err = SetupMetric(autocorrelation)
	if err != nil {
		return nil, errors.Wrap(err, "error setting up Metric")
	}

	return autocorrelation, nil
}

// Subscribe subscribes the Autocorrelation to a Core object.
func (a *Autocorrelation) Subscribe(c *Core) {
	a.correlation.Subscribe(c)
	a.core = c
}

// Config returns the CoreConfig needed.
func (a *Autocorrelation) Config() *CoreConfig {
	return a.correlation.Config()
}

// Push adds a new pair of values for Autocorrelation to consume.
func (a *Autocorrelation) Push(xs ...float64) error {
	if len(xs) != 2 {
		return errors.Errorf(
			"Autocorrelation expected 2 arguments: got %d (%v)",
			len(xs),
			xs,
		)
	}

	x, y := xs[0], xs[1]

	a.core.Lock()
	defer a.core.Unlock()

	if a.queue.Len() >= a.lag {
		tail, err := a.queue.Get()
		if err != nil {
			return errors.Wrap(err, "error popping item from lag queue")
		}

		val := tail.(float64)
		err = a.core.UnsafePush(x, val)
		if err != nil {
			return errors.Wrapf(err, "error pushing (%f, %f) to core", x, val)
		}
	}

	err := a.queue.Put(y)
	if err != nil {
		return errors.Wrapf(err, "error pushing %f to lag queue", y)
	}

	return nil
}

// Value returns the value of the sample autocorrelation.
func (a *Autocorrelation) Value() (float64, error) {
	return a.correlation.Value()
}
