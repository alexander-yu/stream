package moment

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// EWMoment is a metric that tracks the kth exponentially weighted sample central moment.
type EWMoment struct {
	k     int
	decay float64
	core  *Core
}

// NewEWMoment instantiates a EWMoment struct.
func NewEWMoment(k int, decay float64) *EWMoment {
	return &EWMoment{
		k:     k,
		decay: decay,
	}
}

// SetCore sets the Core.
func (m *EWMoment) SetCore(c *Core) {
	m.core = c
}

// IsSetCore returns if the core has been set.
func (m *EWMoment) IsSetCore() bool {
	return m.core != nil
}

// Config returns the CoreConfig needed.
func (m *EWMoment) Config() *CoreConfig {
	return &CoreConfig{
		Sums:  SumsConfig{m.k: true},
		Decay: &m.decay,
	}
}

// String returns a string representation of the metric.
func (m *EWMoment) String() string {
	name := "moment.EWMoment"
	params := []string{
		fmt.Sprintf("k:%v", m.k),
		fmt.Sprintf("decay:%v", m.decay),
	}
	return fmt.Sprintf("%s_{%s}", name, strings.Join(params, ","))
}

// Push adds a new value for EWMoment to consume.
func (m *EWMoment) Push(x float64) error {
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
func (m *EWMoment) Value() (float64, error) {
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
func (m *EWMoment) Clear() {
	if m.IsSetCore() {
		m.core.Clear()
	}
}
