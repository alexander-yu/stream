package joint

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/alexander-yu/stream"
)

// EWMCov is a metric that tracks the sample exponentially weighted covariance.
type EWMCov struct {
	decay float64
	core  *Core
}

// NewEWMCov instantiates a EWMCov struct.
func NewEWMCov(decay float64) *EWMCov {
	return &EWMCov{decay: decay}
}

// SetCore sets the Core.
func (cov *EWMCov) SetCore(c *Core) {
	cov.core = c
}

// IsSetCore returns if the core has been set.
func (cov *EWMCov) IsSetCore() bool {
	return cov.core != nil
}

// Config returns the CoreConfig needed.
func (cov *EWMCov) Config() *CoreConfig {
	return &CoreConfig{
		Sums:   SumsConfig{{1, 1}},
		Window: stream.IntPtr(0),
		Decay:  &cov.decay,
	}
}

// String returns a string representation of the metric.
func (cov *EWMCov) String() string {
	name := "joint.EWMCov"
	return fmt.Sprintf("%s_{decay:%v}", name, cov.decay)
}

// Push adds a new pair of values for EWMCov to consume.
func (cov *EWMCov) Push(xs ...float64) error {
	if !cov.IsSetCore() {
		return errors.New("Core is not set")
	}

	if len(xs) != 2 {
		return errors.Errorf(
			"EWMCov expected 2 arguments: got %d (%v)",
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

// Value returns the value of the sample exponentially weighted covariance.
func (cov *EWMCov) Value() (float64, error) {
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
func (cov *EWMCov) Clear() {
	if cov.IsSetCore() {
		cov.core.Clear()
	}
}
