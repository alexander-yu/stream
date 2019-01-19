package joint

import (
	"fmt"
	"strings"

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
// The lag parameter is of type uint64 due to the internal queue.RingBuffer
// used, whose constructor requires uint64 as input.
func NewAutocorrelation(lag uint64, window int) *Autocorrelation {
	autocorrelation := &Autocorrelation{
		lag:         lag,
		queue:       queue.NewRingBuffer(lag),
		correlation: NewCorrelation(window),
	}

	return autocorrelation
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

// String returns a string representation of the metric.
func (a *Autocorrelation) String() string {
	name := "joint.Autocorrelation"
	params := []string{
		fmt.Sprintf("lag:%v", a.lag),
		fmt.Sprintf("window:%v", a.correlation.window),
	}
	return fmt.Sprintf("%s_{%s}", name, strings.Join(params, ","))
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

	if a.lag == 0 {
		err := a.core.UnsafePush(x, y)
		if err != nil {
			return errors.Wrapf(err, "error pushing (%f, %f) to core", x, y)
		}
		return nil
	}

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

// Clear resets the metric.
func (a *Autocorrelation) Clear() {
	a.correlation.core.Lock()
	defer a.correlation.core.Unlock()
	a.correlation.core.UnsafeClear()
	a.queue.Dispose()
	a.queue = queue.NewRingBuffer(a.lag)
}
