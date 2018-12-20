package moment

import (
	"math"

	"github.com/pkg/errors"
)

// Std is a metric that tracks the sample standard deviation.
type Std struct {
	variance *Moment
}

// NewStd instantiates an Std struct.
func NewStd(window int) (*Std, error) {
	variance, err := NewMoment(2, window)
	if err != nil {
		return nil, errors.Wrap(err, "error creating 2nd Moment")
	}

	return &Std{variance: variance}, nil
}

// Subscribe subscribes the Std to a Core object.
func (s *Std) Subscribe(c *Core) {
	s.variance.Subscribe(c)
}

// Config returns the CoreConfig needed.
func (s *Std) Config() *CoreConfig {
	return s.variance.Config()
}

// Push adds a new value for Std to consume.
func (s *Std) Push(x float64) error {
	err := s.variance.Push(x)
	if err != nil {
		return errors.Wrap(err, "error pushing to core")
	}
	return nil
}

// Value returns the value of the sample standard deviation.
func (s *Std) Value() (float64, error) {
	variance, err := s.variance.Value()
	if err != nil {
		return 0, errors.Wrap(err, "error retrieving 2nd moment")
	}
	return math.Sqrt(variance), nil
}
