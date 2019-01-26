package quantile

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/alexander-yu/stream"
)

// Median keeps track of the median of a stream using order statistic trees.
type Median struct {
	quantile *Quantile
}

// NewMedian instantiates an Median struct. The implementation of the
// underlying order statistic tree can be configured by passing in a constant
// of type Impl.
func NewMedian(window int, impl Impl) (*Median, error) {
	quantile, err := NewQuantile(&Config{
		Window:        stream.IntPtr(window),
		Interpolation: Midpoint.Ptr(),
		Impl:          impl.Ptr(),
	})
	if err != nil {
		return nil, errors.Wrap(err, "error creating Quantile")
	}

	return &Median{quantile: quantile}, nil
}

// String returns a string representation of the metric.
func (m *Median) String() string {
	name := "quantile.Median"
	quantile := fmt.Sprintf("quantile:%v", m.quantile.String())
	return fmt.Sprintf("%s_{%s}", name, quantile)
}

// Push adds a number for calculating the median.
func (m *Median) Push(x float64) error {
	err := m.quantile.Push(x)
	if err != nil {
		return errors.Wrapf(err, "error pushing %f to Quantile", x)
	}
	return nil
}

// Value returns the value of the median.
func (m *Median) Value() (float64, error) {
	value, err := m.quantile.Value(0.5)
	if err != nil {
		return 0, errors.Wrap(err, "error retrieving quantile value")
	}
	return value, nil
}

// Clear resets the metric.
func (m *Median) Clear() {
	m.quantile.Clear()
}
