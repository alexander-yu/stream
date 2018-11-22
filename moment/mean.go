package moment

import (
	"github.com/pkg/errors"

	"github.com/alexander-yu/stream"
)

// Mean is a metric that tracks the mean.
type Mean struct {
	core *stream.Core
}

// Subscribe subscribes the Mean to a Core object.
func (m *Mean) Subscribe(c *stream.Core) {
	m.core = c
}

// Config returns the CoreConfig needed.
func (m *Mean) Config() *stream.CoreConfig {
	return &stream.CoreConfig{
		Sums: stream.SumsConfig{1: true},
	}
}

// Value returns the value of the mean.
func (m *Mean) Value() (float64, error) {
	sum, err := m.core.Sum(1)
	if err != nil {
		return 0, errors.Wrap(err, "error retrieving sum")
	}
	return sum / float64(m.core.Count()), nil
}
