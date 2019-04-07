package moment

import (
	"fmt"
	"math"

	"github.com/pkg/errors"
)

// Kurtosis is a metric that tracks the sample excess kurtosis.
type Kurtosis struct {
	variance *Moment
	moment4  *Moment
	config   *CoreConfig
	core     *Core
}

// NewKurtosis instantiates a Kurtosis struct.
func NewKurtosis(window int) *Kurtosis {
	config := &CoreConfig{
		Sums: SumsConfig{
			2: true,
			4: true,
		},
		Window: &window,
	}

	return &Kurtosis{
		variance: New(2, window),
		moment4:  New(4, window),
		config:   config,
	}
}

// NewGlobalKurtosis instantiates a global Kurtosis struct.
// This is equivalent to calling NewKurtosis(0).
func NewGlobalKurtosis() *Kurtosis {
	return NewKurtosis(0)
}

// SetCore sets the Core.
func (k *Kurtosis) SetCore(c *Core) {
	k.variance.SetCore(c)
	k.moment4.SetCore(c)
	k.core = c
}

// IsSetCore returns if the core has been set.
func (k *Kurtosis) IsSetCore() bool {
	return k.core != nil
}

// Config returns the CoreConfig needed.
func (k *Kurtosis) Config() *CoreConfig {
	return k.config
}

// String returns a string representation of the metric.
func (k *Kurtosis) String() string {
	name := "moment.Kurtosis"
	window := fmt.Sprintf("window:%v", *k.config.Window)
	return fmt.Sprintf("%s_{%s}", name, window)
}

// Push adds a new value for Kurtosis to consume.
func (k *Kurtosis) Push(x float64) error {
	if !k.IsSetCore() {
		return errors.New("Core is not set")
	}

	err := k.core.Push(x)
	if err != nil {
		return errors.Wrap(err, "error pushing to core")
	}
	return nil
}

// Value returns the value of the sample excess kurtosis.
func (k *Kurtosis) Value() (float64, error) {
	if !k.IsSetCore() {
		return 0, errors.New("Core is not set")
	}

	k.core.RLock()
	defer k.core.RUnlock()

	count := float64(k.core.Count())
	if count == 0 {
		return 0, errors.New("no values seen yet")
	}

	variance, err := k.variance.Value()
	if err != nil {
		return 0, errors.Wrap(err, "error retrieving 2nd moment")
	}

	moment, err := k.moment4.Value()
	if err != nil {
		return 0, errors.Wrap(err, "error retrieving 4th moment")
	}

	moment *= (count - 1) / count
	variance *= (count - 1) / count

	return moment/math.Pow(variance, 2) - 3, nil
}

// Clear resets the metric.
func (k *Kurtosis) Clear() {
	if k.IsSetCore() {
		k.core.Clear()
	}
}
