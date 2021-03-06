package quantile

import (
	"fmt"

	"github.com/pkg/errors"
)

// Median keeps track of the median of a stream using order statistics.
type Median struct {
	quantile *Quantile
}

// NewMedian instantiates a Median struct. The implementation of the underlying data
// structure for tracking order statistics can be configured by passing in a constant
// of type Impl.
func NewMedian(window int, options ...Option) (*Median, error) {
	quantile, err := New(window, append(options, InterpolationOption(Midpoint))...)
	if err != nil {
		return nil, errors.Wrap(err, "error creating Quantile")
	}

	return &Median{quantile: quantile}, nil
}

// NewGlobalMedian instantiates a global Median struct.
// This is equivalent to calling NewMedian(0, options...).
func NewGlobalMedian(options ...Option) (*Median, error) {
	return NewMedian(0, options...)
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
