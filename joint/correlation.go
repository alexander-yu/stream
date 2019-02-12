package joint

import (
	"fmt"
	"math"

	"github.com/pkg/errors"
)

// Corr is a metric that tracks the sample Pearson correlation coefficient.
type Corr struct {
	window int
	core   *Core
}

// NewCorr instantiates a Corr struct.
func NewCorr(window int) *Corr {
	return &Corr{window: window}
}

// SetCore sets the Core.
func (corr *Corr) SetCore(c *Core) {
	corr.core = c
}

// IsSetCore returns if the core has been set.
func (corr *Corr) IsSetCore() bool {
	return corr.core != nil
}

// Config returns the CoreConfig needed.
func (corr *Corr) Config() *CoreConfig {
	return &CoreConfig{
		Sums: SumsConfig{
			{1, 1},
			{2, 0},
			{0, 2},
		},
		Window: &corr.window,
	}
}

// String returns a string representation of the metric.
func (corr *Corr) String() string {
	name := "joint.Corr"
	return fmt.Sprintf("%s_{window:%v}", name, corr.window)
}

// Push adds a new pair of values for Corr to consume.
func (corr *Corr) Push(xs ...float64) error {
	if !corr.IsSetCore() {
		return errors.New("Core is not set")
	}

	if len(xs) != 2 {
		return errors.Errorf(
			"Corr expected 2 arguments: got %d (%v)",
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
func (corr *Corr) Value() (float64, error) {
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
func (corr *Corr) Clear() {
	if corr.IsSetCore() {
		corr.core.Clear()
	}
}
