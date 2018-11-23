package moment

import (
	"math"

	"github.com/pkg/errors"

	"github.com/alexander-yu/stream"
)

// Std is a metric that tracks the sample standard deviation.
type Std struct {
	variance *Moment
}

// NewStd creates an Std.
func NewStd() (*Std, error) {
	variance, err := NewMoment(2)
	if err != nil {
		return nil, errors.Wrap(err, "error creating Moment")
	}

	return &Std{variance: variance}, nil
}

// Subscribe subscribes the Std to a Core object.
func (s *Std) Subscribe(c *stream.Core) {
	s.variance.Subscribe(c)
}

// Config returns the CoreConfig needed.
func (s *Std) Config() *stream.CoreConfig {
	return s.variance.Config()
}

// Push is a no-op; Std does not consume values.
func (s *Std) Push(x float64) {}

// Value returns the value of the sample standard deviation.
func (s *Std) Value() (float64, error) {
	variance, err := s.variance.Value()
	if err != nil {
		return 0, errors.Wrap(err, "error retrieving 2nd moment")
	}
	return math.Sqrt(variance), nil
}
