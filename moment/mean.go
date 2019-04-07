package moment

import (
	"fmt"

	"github.com/pkg/errors"
)

// Mean is a metric that tracks the mean.
type Mean struct {
	window int
	core   *Core
}

// NewMean instantiates a Mean struct.
func NewMean(window int) *Mean {
	return &Mean{window: window}
}

// NewGlobalMean instantiates a global Mean struct.
// This is equivalent to calling NewMean(0).
func NewGlobalMean() *Mean {
	return NewMean(0)
}

// SetCore sets the Core.
func (m *Mean) SetCore(c *Core) {
	m.core = c
}

// IsSetCore returns if the core has been set.
func (m *Mean) IsSetCore() bool {
	return m.core != nil
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
	if !m.IsSetCore() {
		return errors.New("Core is not set")
	}

	err := m.core.Push(x)
	if err != nil {
		return errors.Wrap(err, "error pushing to core")
	}
	return nil
}

// Value returns the value of the mean.
func (m *Mean) Value() (float64, error) {
	if !m.IsSetCore() {
		return 0, errors.New("Core is not set")
	}

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
	if m.IsSetCore() {
		m.core.Clear()
	}
}
