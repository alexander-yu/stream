package moment

import (
	"github.com/pkg/errors"

	"stream"
)

// Mean is a metric that tracks the mean.
type Mean struct {
	sum  *Moment
	core *stream.Core
}

// NewMean creates a Mean.
func NewMean() (*Mean, error) {
	sum, err := NewMoment(1)
	if err != nil {
		return nil, errors.Wrap(err, "error creating Moment")
	}

	return &Mean{sum: sum}, nil
}

// Subscribe subscribes the Mean to a Core object.
func (m *Mean) Subscribe(c *stream.Core) {
	m.sum.Subscribe(c)
	m.core = c
}

// Config returns the CoreConfig needed.
func (m *Mean) Config() *stream.CoreConfig {
	return m.sum.Config()
}

// Value returns the value of the mean.
func (m *Mean) Value() (float64, error) {
	sum, err := m.sum.Value()
	if err != nil {
		return 0, errors.Wrap(err, "error retrieving 1st moment")
	}
	return sum / float64(m.core.Count()), nil
}
