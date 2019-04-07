package joint

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/workiva/go-datastructures/queue"
)

// Autocorr is a metric that tracks the sample autocorrelation.
// It does not satisfy the JointMetric interface, but rather the univariate
// Metric interface (SimpleMetric in particular), since it only tracks a single
// variable.
type Autocorr struct {
	lag   int
	queue *queue.RingBuffer
	corr  *Corr
	core  *Core
}

// NewAutocorr instantiates an Autocorr struct.
func NewAutocorr(lag int, window int) (*Autocorr, error) {
	if lag < 0 {
		return nil, errors.Errorf("%d is a negative lag", lag)
	}

	return &Autocorr{
		lag:   lag,
		queue: queue.NewRingBuffer(uint64(lag)),
		corr:  NewCorr(window),
	}, nil
}

// NewGlobalAutocorr instantiates a global Autocorr struct.
// This is equivalent to calling NewAutocorr(lag, 0).
func NewGlobalAutocorr(lag int) (*Autocorr, error) {
	return NewAutocorr(lag, 0)
}

// SetCore sets the Core.
func (a *Autocorr) SetCore(c *Core) {
	a.corr.SetCore(c)
	a.core = c
}

// IsSetCore returns if the core has been set.
func (a *Autocorr) IsSetCore() bool {
	return a.core != nil
}

// Config returns the CoreConfig needed.
func (a *Autocorr) Config() *CoreConfig {
	return a.corr.Config()
}

// String returns a string representation of the metric.
func (a *Autocorr) String() string {
	name := "joint.Autocorr"
	params := []string{
		fmt.Sprintf("lag:%v", a.lag),
		fmt.Sprintf("window:%v", a.corr.window),
	}
	return fmt.Sprintf("%s_{%s}", name, strings.Join(params, ","))
}

// Push adds a new value for Autocorr to consume.
// Autocorr only takes one value, because we're calculating
// the lagged correlation of a series of data against itself.
func (a *Autocorr) Push(x float64) error {
	if !a.IsSetCore() {
		return errors.New("Core is not set")
	}

	a.core.Lock()
	defer a.core.Unlock()

	if a.lag == 0 {
		err := a.core.UnsafePush(x, x)
		if err != nil {
			return errors.Wrap(err, "error pushing to core")
		}
		return nil
	}

	if a.queue.Len() >= uint64(a.lag) {
		tail, err := a.queue.Get()
		if err != nil {
			return errors.Wrap(err, "error popping item from lag queue")
		}

		val := tail.(float64)
		err = a.core.UnsafePush(x, val)
		if err != nil {
			return errors.Wrap(err, "error pushing to core")
		}
	}

	err := a.queue.Put(x)
	if err != nil {
		return errors.Wrapf(err, "error pushing %f to lag queue", x)
	}

	return nil
}

// Value returns the value of the sample autocorrelation.
func (a *Autocorr) Value() (float64, error) {
	if !a.IsSetCore() {
		return 0, errors.New("Core is not set")
	} else if a.corr.core.Count() == 0 {
		return 0, errors.New(fmt.Sprintf(
			"Not enough values seen; at least %d observations must be made",
			a.lag+1,
		))
	}

	return a.corr.Value()
}

// Clear resets the metric.
func (a *Autocorr) Clear() {
	if a.IsSetCore() {
		a.corr.core.Lock()
		defer a.corr.core.Unlock()
		a.corr.core.UnsafeClear()
		a.queue.Dispose()
		a.queue = queue.NewRingBuffer(uint64(a.lag))
	}
}
