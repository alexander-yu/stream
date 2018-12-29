package joint

import (
	"math"

	"github.com/pkg/errors"
)

// Correlation is a metric that tracks the sample Pearson correlation coefficient.
type Correlation struct {
	window int
	core   *Core
}

// NewCorrelation instantiates a Correlation struct.
func NewCorrelation(window int) (*Correlation, error) {
	correlation := &Correlation{window: window}

	err := SetupMetric(correlation)
	if err != nil {
		return nil, errors.Wrap(err, "error setting up Metric")
	}

	return correlation, nil
}

// Subscribe subscribes the Correlation to a Core object.
func (cov *Correlation) Subscribe(c *Core) {
	cov.core = c
}

// Config returns the CoreConfig needed.
func (cov *Correlation) Config() *CoreConfig {
	return &CoreConfig{
		Sums: SumsConfig{
			{1, 1},
			{2, 0},
			{0, 2},
		},
		Window: &cov.window,
	}
}

// Push adds a new pair of values for Correlation to consume.
func (cov *Correlation) Push(xs ...float64) error {
	if len(xs) != 2 {
		return errors.Errorf(
			"Correlation expected 2 arguments: got %d (%v)",
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

// Value returns the value of the sample Pearson correlation coefficient.
func (cov *Correlation) Value() (float64, error) {
	cov.core.RLock()
	defer cov.core.RUnlock()

	covariance, err := cov.core.Sum(1, 1)
	if err != nil {
		return 0, errors.Wrap(err, "error retrieving sum for {1, 1}")
	}

	xVar, err := cov.core.Sum(2, 0)
	if err != nil {
		return 0, errors.Wrap(err, "error retrieving sum for {2, 0}")
	}

	yVar, err := cov.core.Sum(0, 2)
	if err != nil {
		return 0, errors.Wrap(err, "error retrieving sum for {0, 2}")
	}

	return covariance / math.Sqrt(xVar*yVar), nil
}
