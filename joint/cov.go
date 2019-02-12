package joint

import (
	"fmt"

	"github.com/pkg/errors"
)

// Cov is a metric that tracks the sample covariance.
type Cov struct {
	window int
	core   *Core
}

// NewCov instantiates a Cov struct.
func NewCov(window int) *Cov {
	return &Cov{window: window}
}

// SetCore sets the Core.
func (cov *Cov) SetCore(c *Core) {
	cov.core = c
}

// IsSetCore returns if the core has been set.
func (cov *Cov) IsSetCore() bool {
	return cov.core != nil
}

// Config returns the CoreConfig needed.
func (cov *Cov) Config() *CoreConfig {
	return &CoreConfig{
		Sums:   SumsConfig{{1, 1}},
		Window: &cov.window,
	}
}

// String returns a string representation of the metric.
func (cov *Cov) String() string {
	name := "joint.Cov"
	return fmt.Sprintf("%s_{window:%v}", name, cov.window)
}

// Push adds a new pair of values for Cov to consume.
func (cov *Cov) Push(xs ...float64) error {
	if !cov.IsSetCore() {
		return errors.New("Core is not set")
	}

	if len(xs) != 2 {
		return errors.Errorf(
			"Cov expected 2 arguments: got %d (%v)",
			len(xs),
			xs,
		)
	}

	err := cov.core.Push(xs...)
	if err != nil {
		return errors.Wrap(err, "error pushing to core")
	}
	return nil
}

// Value returns the value of the sample covariance.
func (cov *Cov) Value() (float64, error) {
	if !cov.IsSetCore() {
		return 0, errors.New("Core is not set")
	}

	cov.core.RLock()
	defer cov.core.RUnlock()

	covariance, err := cov.core.Sum(1, 1)
	if err != nil {
		return 0, errors.Wrap(err, "error retrieving sum")
	}

	count := cov.core.Count()
	covariance /= (float64(count) - 1.)

	return covariance, nil
}

// Clear resets the metric.
func (cov *Cov) Clear() {
	if cov.IsSetCore() {
		cov.core.Clear()
	}
}
