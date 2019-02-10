package moment

import (
	"fmt"
	"math"

	"github.com/pkg/errors"
)

// EWMStd is a metric that tracks the exponentially weighted sample standard deviation.
type EWMStd struct {
	variance *EWMMoment
}

// NewEWMStd instantiates an EWMStd struct.
func NewEWMStd(decay float64) *EWMStd {
	return &EWMStd{variance: NewEWMMoment(2, decay)}
}

// SetCore sets the Core.
func (s *EWMStd) SetCore(c *Core) {
	s.variance.SetCore(c)
}

// IsSetCore returns if the core has been set.
func (s *EWMStd) IsSetCore() bool {
	return s.variance.IsSetCore()
}

// Config returns the CoreConfig needed.
func (s *EWMStd) Config() *CoreConfig {
	return s.variance.Config()
}

// String returns a string representation of the metric.
func (s *EWMStd) String() string {
	name := "moment.EWMStd"
	decay := fmt.Sprintf("decay:%v", *s.variance.Config().Decay)
	return fmt.Sprintf("%s_{%s}", name, decay)
}

// Push adds a new value for EWMStd to consume.
func (s *EWMStd) Push(x float64) error {
	if !s.IsSetCore() {
		return errors.New("Core is not set")
	}

	err := s.variance.Push(x)
	if err != nil {
		return errors.Wrap(err, "error pushing to core")
	}
	return nil
}

// Value returns the value of the exponentially weighted sample standard deviation.
func (s *EWMStd) Value() (float64, error) {
	if !s.IsSetCore() {
		return 0, errors.New("Core is not set")
	}

	variance, err := s.variance.Value()
	if err != nil {
		return 0, errors.Wrap(err, "error retrieving 2nd moment")
	}
	return math.Sqrt(variance), nil
}

// Clear resets the metric.
func (s *EWMStd) Clear() {
	if s.IsSetCore() {
		s.variance.Clear()
	}
}
