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
	variance := &Moment{K: 2, Window: window}
	moment4 := &Moment{K: 4, Window: window}
	config := &CoreConfig{
		Sums: SumsConfig{
			2: true,
			4: true,
		},
		Window: &window,
	}

	return &Kurtosis{
		variance: variance,
		moment4:  moment4,
		config:   config,
	}
}

// Subscribe subscribes the Kurtosis to a Core object.
func (k *Kurtosis) Subscribe(c *Core) {
	k.variance.Subscribe(c)
	k.moment4.Subscribe(c)
	k.core = c
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
	err := k.core.Push(x)
	if err != nil {
		return errors.Wrap(err, "error pushing to core")
	}
	return nil
}

// Value returns the value of the sample excess kurtosis.
func (k *Kurtosis) Value() (float64, error) {
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
	k.core.Clear()
}
