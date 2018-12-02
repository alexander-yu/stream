package moment

import (
	"math"

	"github.com/pkg/errors"

	"github.com/alexander-yu/stream"
)

// Kurtosis is a metric that tracks the kurtosis.
type Kurtosis struct {
	variance *Moment
	moment4  *Moment
	config   *stream.CoreConfig
	core     *stream.Core
}

// NewKurtosis instantiates a Kurtosis struct.
func NewKurtosis(window int) (*Kurtosis, error) {
	variance, err := NewMoment(2, window)
	if err != nil {
		return nil, errors.Wrap(err, "error creating 2nd Moment")
	}

	moment4, err := NewMoment(4, window)
	if err != nil {
		return nil, errors.Wrap(err, "error creating 4th Moment")
	}

	config, err := stream.MergeConfigs(variance.Config(), moment4.Config())
	if err != nil {
		return nil, errors.Wrap(err, "error merging configs")
	}

	return &Kurtosis{
		variance: variance,
		moment4:  moment4,
		config:   config,
	}, nil
}

// Subscribe subscribes the Kurtosis to a Core object.
func (k *Kurtosis) Subscribe(c *stream.Core) {
	k.variance.Subscribe(c)
	k.moment4.Subscribe(c)
	k.core = c
}

// Config returns the CoreConfig needed.
func (k *Kurtosis) Config() *stream.CoreConfig {
	return k.config
}

// Push is a no-op; Kurtosis does not consume values.
func (k *Kurtosis) Push(x float64) error { return nil }

// Value returns the value of the sample excess kurtosis.
func (k *Kurtosis) Value() (float64, error) {
	count := float64(k.core.Count())
	if count == 0 {
		return 0, errors.New("no values seen yet")
	}

	variance, err := k.variance.Value()
	if err != nil {
		return 0, errors.Wrap(err, "error retrieving 2nd moment")
	}

	moment, err := k.moment4.Value()
	if err != nil {
		return 0, errors.Wrap(err, "error retrieving 4th moment")
	}

	moment *= (count - 1) / count
	variance *= (count - 1) / count

	return moment/math.Pow(variance, 2) - 3, nil
}
