package joint

import (
	"fmt"
	"math"

	"github.com/pkg/errors"

	"github.com/alexander-yu/stream"
)

// EWMCorr is a metric that tracks the sample Pearson correlation coefficient.
type EWMCorr struct {
	decay float64
	core  *Core
}

// NewEWMCorr instantiates a EWMCorr struct.
func NewEWMCorr(decay float64) *EWMCorr {
	return &EWMCorr{decay: decay}
}

// SetCore sets the Core.
func (corr *EWMCorr) SetCore(c *Core) {
	corr.core = c
}

// IsSetCore returns if the core has been set.
func (corr *EWMCorr) IsSetCore() bool {
	return corr.core != nil
}

// Config returns the CoreConfig needed.
func (corr *EWMCorr) Config() *CoreConfig {
	return &CoreConfig{
		Sums: SumsConfig{
			{1, 1},
			{2, 0},
			{0, 2},
		},
		Window: stream.IntPtr(0),
		Decay:  &corr.decay,
	}
}

// String returns a string representation of the metric.
func (corr *EWMCorr) String() string {
	name := "joint.EWMCorr"
	return fmt.Sprintf("%s_{decay:%v}", name, corr.decay)
}

// Push adds a new pair of values for EWMCorr to consume.
func (corr *EWMCorr) Push(xs ...float64) error {
	if !corr.IsSetCore() {
		return errors.New("Core is not set")
	}

	if len(xs) != 2 {
		return errors.Errorf(
			"EWMCorr expected 2 arguments: got %d (%v)",
			len(xs),
			xs,
		)
	}

	err := corr.core.Push(xs...)
	if err != nil {
		return errors.Wrap(err, "error pushing to core")
	}
	return nil
}

// Value returns the value of the sample Pearson correlation coefficient.
func (corr *EWMCorr) Value() (float64, error) {
	if !corr.IsSetCore() {
		return 0, errors.New("Core is not set")
	}

	corr.core.RLock()
	defer corr.core.RUnlock()

	// this is technically not the covariance, as it is not normalized by
	// the sample size (minus 1), but the denominator is cancelled out
	// when dividing by the sqrt of the variances, so we can avoid extra
	// float ops here
	cov, err := corr.core.Sum(1, 1)
	if err != nil {
		return 0, errors.Wrap(err, "error retrieving sum for {1, 1}")
	}

	// ditto with the "variance" variables here, as with above
	xVar, err := corr.core.Sum(2, 0)
	if err != nil {
		return 0, errors.Wrap(err, "error retrieving sum for {2, 0}")
	}

	yVar, err := corr.core.Sum(0, 2)
	if err != nil {
		return 0, errors.Wrap(err, "error retrieving sum for {0, 2}")
	}

	return cov / math.Sqrt(xVar*yVar), nil
}

// Clear resets the metric.
func (corr *EWMCorr) Clear() {
	if corr.IsSetCore() {
		corr.core.Clear()
	}
}
