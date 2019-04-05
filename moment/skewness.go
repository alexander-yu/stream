package moment

import (
	"fmt"
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
func NewSkewness(window int) *Skewness {
	config := &CoreConfig{
		Sums: SumsConfig{
			2: true,
			3: true,
		},
		Window: &window,
	}

	return &Skewness{
		variance: New(2, window),
		moment3:  New(3, window),
		config:   config,
	}
}

// SetCore sets the Core.
func (s *Skewness) SetCore(c *Core) {
	s.variance.SetCore(c)
	s.moment3.SetCore(c)
	s.core = c
}

// IsSetCore returns if the core has been set.
func (s *Skewness) IsSetCore() bool {
	return s.core != nil
}

// Config returns the CoreConfig needed.
func (s *Skewness) Config() *CoreConfig {
	return s.config
}

// String returns a string representation of the metric.
func (s *Skewness) String() string {
	name := "moment.Skewness"
	window := fmt.Sprintf("window:%v", *s.config.Window)
	return fmt.Sprintf("%s_{%s}", name, window)
}

// Push adds a new value for Skewness to consume.
func (s *Skewness) Push(x float64) error {
	if !s.IsSetCore() {
		return errors.New("Core is not set")
	}

	err := s.core.Push(x)
	if err != nil {
		return errors.Wrap(err, "error pushing to core")
	}
	return nil
}

// Value returns the value of the adjusted Fisher-Pearson sample skewness.
func (s *Skewness) Value() (float64, error) {
	if !s.IsSetCore() {
		return 0, errors.New("Core is not set")
	}

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

// Clear resets the metric.
func (s *Skewness) Clear() {
	if s.IsSetCore() {
		s.core.Clear()
	}
}
