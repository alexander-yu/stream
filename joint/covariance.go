package joint

import (
	"fmt"

	"github.com/pkg/errors"
)

// Covariance is a metric that tracks the sample covariance.
type Covariance struct {
	window int
	core   *Core
}

// NewCovariance instantiates a Covariance struct.
func NewCovariance(window int) *Covariance {
	return &Covariance{window: window}
}

// SetCore sets the Core.
func (cov *Covariance) SetCore(c *Core) {
	cov.core = c
}

// IsSetCore returns if the core has been set.
func (cov *Covariance) IsSetCore() bool {
	return cov.core != nil
}

// Config returns the CoreConfig needed.
func (cov *Covariance) Config() *CoreConfig {
	return &CoreConfig{
		Sums:   SumsConfig{{1, 1}},
		Window: &cov.window,
	}
}

// String returns a string representation of the metric.
func (cov *Covariance) String() string {
	name := "joint.Covariance"
	return fmt.Sprintf("%s_{window:%v}", name, cov.window)
}

// Push adds a new pair of values for Covariance to consume.
func (cov *Covariance) Push(xs ...float64) error {
	if !cov.IsSetCore() {
		return errors.New("Core is not set")
	}

	if len(xs) != 2 {
		return errors.Errorf(
			"Covariance expected 2 arguments: got %d (%v)",
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
func (cov *Covariance) Value() (float64, error) {
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
func (cov *Covariance) Clear() {
	if cov.IsSetCore() {
		cov.core.Clear()
	}
}
