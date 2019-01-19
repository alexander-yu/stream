package moment

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// Moment is a metric that tracks the kth sample central moment.
type Moment struct {
	K      int
	Window int
	core   *Core
}

// NewMoment instantiates a Moment struct.
func NewMoment(k int, window uint64) *Moment {
	return &Moment{
		k:      k,
		window: window,
	}
}

// Subscribe subscribes the Moment to a Core object.
func (m *Moment) Subscribe(c *Core) {
	m.core = c
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
	err := m.core.Push(x)
	if err != nil {
		return errors.Wrap(err, "error pushing to core")
	}
	return nil
}

// Value returns the value of the kth sample central moment.
func (m *Moment) Value() (float64, error) {
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
	m.core.Clear()
}
