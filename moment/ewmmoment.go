package moment

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// EWMMoment is a metric that tracks the kth exponentially weighted sample central moment.
type EWMMoment struct {
	k     int
	decay float64
	core  *Core
}

// NewEWMMoment instantiates a EWMMoment struct.
func NewEWMMoment(k int, decay float64) *EWMMoment {
	return &EWMMoment{
		k:     k,
		decay: decay,
	}
}

// SetCore sets the Core.
func (m *EWMMoment) SetCore(c *Core) {
	m.core = c
}

// IsSetCore returns if the core has been set.
func (m *EWMMoment) IsSetCore() bool {
	return m.core != nil
}

// Config returns the CoreConfig needed.
func (m *EWMMoment) Config() *CoreConfig {
	return &CoreConfig{
		Sums:  SumsConfig{m.k: true},
		Decay: &m.decay,
	}
}

// String returns a string representation of the metric.
func (m *EWMMoment) String() string {
	name := "moment.EWMMoment"
	params := []string{
		fmt.Sprintf("k:%v", m.k),
		fmt.Sprintf("decay:%v", m.decay),
	}
	return fmt.Sprintf("%s_{%s}", name, strings.Join(params, ","))
}

// Push adds a new value for EWMMoment to consume.
func (m *EWMMoment) Push(x float64) error {
	if !m.IsSetCore() {
		return errors.New("Core is not set")
	}

	err := m.core.Push(x)
	if err != nil {
		return errors.Wrap(err, "error pushing to core")
	}
	return nil
}

// Value returns the value of the kth exponentially weighted sample central moment.
func (m *EWMMoment) Value() (float64, error) {
	if !m.IsSetCore() {
		return 0, errors.New("Core is not set")
	}

	m.core.RLock()
	defer m.core.RUnlock()

	moment, err := m.core.Sum(m.k)
	if err != nil {
		return 0, errors.Wrap(err, "error retrieving sum")
	}

	count := m.core.Count()
	moment /= (float64(count) - 1.)

	return moment, nil
}

// Clear resets the metric.
func (m *EWMMoment) Clear() {
	if m.IsSetCore() {
		m.core.Clear()
	}
}
