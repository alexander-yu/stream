package moment

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// Moment is a metric that tracks the kth sample central moment.
type Moment struct {
	k      int
	window int
	core   *Core
}

// New instantiates a Moment struct.
func New(k int, window int) *Moment {
	return &Moment{
		k:      k,
		window: window,
	}
}

// NewGlobal instantiates a global Moment struct.
// This is equivalent to calling New(k, 0).
func NewGlobal(k int) *Moment {
	return New(k, 0)
}

// SetCore sets the Core.
func (m *Moment) SetCore(c *Core) {
	m.core = c
}

// IsSetCore returns if the core has been set.
func (m *Moment) IsSetCore() bool {
	return m.core != nil
}

// Config returns the CoreConfig needed.
func (m *Moment) Config() *CoreConfig {
	return &CoreConfig{
		Sums:   SumsConfig{m.k: true},
		Window: &m.window,
	}
}

// String returns a string representation of the metric.
func (m *Moment) String() string {
	name := "moment.Moment"
	params := []string{
		fmt.Sprintf("k:%v", m.k),
		fmt.Sprintf("window:%v", m.window),
	}
	return fmt.Sprintf("%s_{%s}", name, strings.Join(params, ","))
}

// Push adds a new value for Moment to consume.
func (m *Moment) Push(x float64) error {
	if !m.IsSetCore() {
		return errors.New("Core is not set")
	}

	err := m.core.Push(x)
	if err != nil {
		return errors.Wrap(err, "error pushing to core")
	}
	return nil
}

// Value returns the value of the kth sample central moment.
func (m *Moment) Value() (float64, error) {
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
func (m *Moment) Clear() {
	if m.IsSetCore() {
		m.core.Clear()
	}
}
