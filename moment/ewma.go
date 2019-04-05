package moment

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/alexander-yu/stream"
)

// EWMA is a metric that tracks the exponentially weighted moving average.
type EWMA struct {
	decay float64
	core  *Core
}

// NewEWMA instantiates a EWMA struct.
func NewEWMA(decay float64) *EWMA {
	return &EWMA{decay: decay}
}

// SetCore sets the Core.
func (a *EWMA) SetCore(c *Core) {
	a.core = c
}

// IsSetCore returns if the core has been set.
func (a *EWMA) IsSetCore() bool {
	return a.core != nil
}

// Config returns the CoreConfig needed.
func (a *EWMA) Config() *CoreConfig {
	return &CoreConfig{
		Window: stream.IntPtr(0),
		Decay:  &a.decay,
	}
}

// String returns a string representation of the metric.
func (a *EWMA) String() string {
	name := "moment.EWMA"
	decay := fmt.Sprintf("decay:%v", a.decay)
	return fmt.Sprintf("%s_{%s}", name, decay)
}

// Push adds a new value for EWMA to consume.
func (a *EWMA) Push(x float64) error {
	if !a.IsSetCore() {
		return errors.New("Core is not set")
	}

	err := a.core.Push(x)
	if err != nil {
		return errors.Wrap(err, "error pushing to core")
	}
	return nil
}

// Value returns the value of the exponentially weighted moving average.
func (a *EWMA) Value() (float64, error) {
	if !a.IsSetCore() {
		return 0, errors.New("Core is not set")
	}

	a.core.RLock()
	defer a.core.RUnlock()

	ewma, err := a.core.Mean()
	if err != nil {
		return 0, errors.Wrap(err, "error retrieving sum")
	}
	return ewma, nil
}

// Clear resets the metric.
func (a *EWMA) Clear() {
	if a.IsSetCore() {
		a.core.Clear()
	}
}
