package moment

import (
	"github.com/pkg/errors"

	"github.com/alexander-yu/stream"
)

// Mean is a metric that tracks the mean.
type Mean struct {
	window int
	core   *stream.Core
}

// NewMean instantiates a Mean struct.
func NewMean(window int) (*Mean, error) {
	if window < 0 {
		return nil, errors.Errorf("%d is a negative window", window)
	}
	return &Mean{window: window}, nil
}

// Subscribe subscribes the Mean to a Core object.
func (m *Mean) Subscribe(c *stream.Core) {
	m.core = c
}

// Config returns the CoreConfig needed.
func (m *Mean) Config() *stream.CoreConfig {
	return &stream.CoreConfig{
		Sums:   stream.SumsConfig{1: true},
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
