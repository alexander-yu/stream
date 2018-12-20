package moment

import (
	"github.com/pkg/errors"
)

// Mean is a metric that tracks the mean.
type Mean struct {
	window int
	core   *Core
}

// NewMean instantiates a Mean struct.
func NewMean(window int) (*Mean, error) {
	mean := &Mean{window: window}
	err := SetupMetric(mean)
	if err != nil {
		return nil, errors.Wrap(err, "error setting up Metric")
	}
	return mean, nil
}

// Subscribe subscribes the Mean to a Core object.
func (m *Mean) Subscribe(c *Core) {
	m.core = c
}

// Config returns the CoreConfig needed.
func (m *Mean) Config() *CoreConfig {
	return &CoreConfig{
		Sums:   SumsConfig{1: true},
		Window: &m.window,
	}
}

// Push adds a new value for Mean to consume.
func (m *Mean) Push(x float64) error {
	err := m.core.Push(x)
	if err != nil {
		return errors.Wrap(err, "error pushing to core")
	}
	return nil
}

// Value returns the value of the mean.
func (m *Mean) Value() (float64, error) {
	m.core.RLock()
	defer m.core.RUnlock()

	count := m.core.Count()
	if count == 0 {
		return 0, errors.New("no values seen yet")
	}

	sum, err := m.core.Sum(1)
	if err != nil {
		return 0, errors.Wrap(err, "error retrieving sum")
	}
	return sum / float64(count), nil
}
