package moment

import (
	"fmt"

	"github.com/pkg/errors"
)

// Mean is a metric that tracks the mean.
type Mean struct {
	Window int
	core   *Core
}

// NewMean instantiates a Mean struct.
func NewMean(window uint64) *Mean {
	return &Mean{window: window}
}

// Subscribe subscribes the Mean to a Core object.
func (m *Mean) Subscribe(c *Core) {
	m.core = c
}

// Config returns the CoreConfig needed.
func (m *Mean) Config() *CoreConfig {
	return &CoreConfig{
		Window: &m.window,
	}
}

// String returns a string representation of the metric.
func (m *Mean) String() string {
	name := "moment.Mean"
	window := fmt.Sprintf("window:%v", m.window)
	return fmt.Sprintf("%s_{%s}", name, window)
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

	mean, err := m.core.Mean()
	if err != nil {
		return 0, errors.Wrap(err, "error retrieving sum")
	}
	return mean, nil
}

// Clear resets the metric.
func (m *Mean) Clear() {
	m.core.Clear()
}
