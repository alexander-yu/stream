package quantile

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/alexander-yu/stream"
	"github.com/alexander-yu/stream/quantile/ost"
)

// OSTMedian keeps track of the median of a stream using order statistic trees.
type OSTMedian struct {
	quantile *OSTQuantile
}

// NewOSTMedian instantiates an OSTMedian struct. The implementation of the
// underlying order statistic tree can be configured by passing in a constant
// of type ost.Impl.
func NewOSTMedian(window int, impl ost.Impl) (*OSTMedian, error) {
	quantile, err := NewOSTQuantile(&Config{
		Quantile:      stream.FloatPtr(0.5),
		Window:        stream.IntPtr(window),
		Interpolation: Midpoint.Ptr(),
	}, impl)
	if err != nil {
		return nil, errors.Wrap(err, "error creating OSTQuantile")
	}

	return &OSTMedian{quantile: quantile}, nil
}

// String returns a string representation of the metric.
func (m *OSTMedian) String() string {
	name := "quantile.OSTMedian"
	quantile := fmt.Sprintf("quantile:%v", m.quantile.String())
	return fmt.Sprintf("%s_{%s}", name, quantile)
}

// Push adds a number for calculating the median.
func (m *OSTMedian) Push(x float64) error {
	err := m.quantile.Push(x)
	if err != nil {
		return errors.Wrapf(err, "error pushing %f to OSTQuantile", x)
	}
	return nil
}

// Value returns the value of the median.
func (m *OSTMedian) Value() (float64, error) {
	value, err := m.quantile.Value()
	if err != nil {
		return 0, errors.Wrap(err, "error retrieving quantile value")
	}
	return value, nil
}
