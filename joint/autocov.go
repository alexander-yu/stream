package joint

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/workiva/go-datastructures/queue"
)

// Autocov is a metric that tracks the sample autocovariance.
// It does not satisfy the JointMetric interface, but rather the univariate
// Metric interface (SimpleMetric in particular), since it only tracks a single
// variable.
type Autocov struct {
	lag   int
	queue *queue.RingBuffer
	cov   *Cov
	core  *Core
}

// NewAutocov instantiates an Autocov struct.
func NewAutocov(lag int, window int) (*Autocov, error) {
	if lag < 0 {
		return nil, errors.Errorf("%d is a negative lag", lag)
	}

	return &Autocov{
		lag:   lag,
		queue: queue.NewRingBuffer(uint64(lag)),
		cov:   NewCov(window),
	}, nil
}

// NewGlobalAutocov instantiates a global Autocov struct.
// This is equivalent to calling NewAutocov(lag, 0).
func NewGlobalAutocov(lag int) (*Autocov, error) {
	return NewAutocov(lag, 0)
}

// SetCore sets the Core.
func (a *Autocov) SetCore(c *Core) {
	a.cov.SetCore(c)
	a.core = c
}

// IsSetCore returns if the core has been set.
func (a *Autocov) IsSetCore() bool {
	return a.core != nil
}

// Config returns the CoreConfig needed.
func (a *Autocov) Config() *CoreConfig {
	return a.cov.Config()
}

// String returns a string representation of the metric.
func (a *Autocov) String() string {
	name := "joint.Autocov"
	params := []string{
		fmt.Sprintf("lag:%v", a.lag),
		fmt.Sprintf("window:%v", a.cov.window),
	}
	return fmt.Sprintf("%s_{%s}", name, strings.Join(params, ","))
}

// Push adds a new value for Autocov to consume.
// Autocov only takes one value, because we're calculating
// the lagged covelation of a series of data against itself.
func (a *Autocov) Push(x float64) error {
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

// Value returns the value of the sample autocovelation.
func (a *Autocov) Value() (float64, error) {
	if !a.IsSetCore() {
		return 0, errors.New("Core is not set")
	} else if a.cov.core.Count() == 0 {
		return 0, errors.Errorf(
			"Not enough values seen; at least %d observations must be made",
			a.lag+1,
		)
	}

	return a.cov.Value()
}

// Clear resets the metric.
func (a *Autocov) Clear() {
	if a.IsSetCore() {
		a.cov.core.Lock()
		defer a.cov.core.Unlock()
		a.cov.core.UnsafeClear()
		a.queue.Dispose()
		a.queue = queue.NewRingBuffer(uint64(a.lag))
	}
}
