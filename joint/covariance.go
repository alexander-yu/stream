package joint

import (
	"fmt"

	"github.com/pkg/errors"
)

// Covariance is a metric that tracks the sample covariance.
type Covariance struct {
	Window int
	core   *Core
}

// Subscribe subscribes the Covariance to a Core object.
func (cov *Covariance) Subscribe(c *Core) {
	cov.core = c
}

// Config returns the CoreConfig needed.
func (cov *Covariance) Config() *CoreConfig {
	return &CoreConfig{
		Sums:   SumsConfig{{1, 1}},
		Window: &cov.Window,
	}
}

// String returns a string representation of the metric.
func (cov *Covariance) String() string {
	name := "joint.Covariance"
	return fmt.Sprintf("%s_{window:%v}", name, cov.Window)
}

// Push adds a new pair of values for Covariance to consume.
func (cov *Covariance) Push(xs ...float64) error {
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
	cov.core.Clear()
}
