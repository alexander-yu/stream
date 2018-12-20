package moment

import (
	"math"

	"github.com/pkg/errors"
)

// Skewness is a metric that tracks the adjusted Fisher-Pearson sample skewness.
type Skewness struct {
	variance *Moment
	moment3  *Moment
	config   *CoreConfig
	core     *Core
}

// NewSkewness instantiates a Skewness struct.
func NewSkewness(window int) (*Skewness, error) {
	variance, err := NewMoment(2, window)
	if err != nil {
		return nil, errors.Wrap(err, "error creating 2nd Moment")
	}

	moment3, err := NewMoment(3, window)
	if err != nil {
		return nil, errors.Wrap(err, "error creating 3rd Moment")
	}

	config, err := MergeConfigs(variance.Config(), moment3.Config())
	if err != nil {
		return nil, errors.Wrap(err, "error merging configs")
	}

	return &Skewness{
		variance: variance,
		moment3:  moment3,
		config:   config,
	}, nil
}

// Subscribe subscribes the Skewness to a Core object.
func (s *Skewness) Subscribe(c *Core) {
	s.variance.Subscribe(c)
	s.moment3.Subscribe(c)
	s.core = c
}

// Config returns the CoreConfig needed.
func (s *Skewness) Config() *CoreConfig {
	return s.config
}

// Push adds a new value for Skewness to consume.
func (s *Skewness) Push(x float64) error {
	err := s.core.Push(x)
	if err != nil {
		return errors.Wrap(err, "error pushing to core")
	}
	return nil
}

// Value returns the value of the adjusted Fisher-Pearson sample skewness.
func (s *Skewness) Value() (float64, error) {
	s.core.RLock()
	defer s.core.RUnlock()

	count := float64(s.core.Count())
	if count == 0 {
		return 0, errors.New("no values seen yet")
	}

	variance, err := s.variance.Value()
	if err != nil {
		return 0, errors.Wrap(err, "error retrieving 2nd moment")
	}
	variance *= (count - 1) / count

	moment, err := s.moment3.Value()
	if err != nil {
		return 0, errors.Wrap(err, "error retrieving 3rd moment")
	}
	moment *= (count - 1) / count

	adjust := math.Sqrt(count*(count-1)) / (count - 2)
	return adjust * moment / math.Pow(variance, 1.5), nil
}
