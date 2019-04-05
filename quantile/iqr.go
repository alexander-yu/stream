package quantile

import (
	"fmt"

	"github.com/pkg/errors"
)

// IQR keeps track of the interquartile range of a stream using order statistics.
type IQR struct {
	quantile *Quantile
}

// NewIQR instantiates an IQR struct.
func NewIQR(window int, options ...Option) (*IQR, error) {
	quantile, err := NewQuantile(window, append(options, InterpolationOption(Midpoint))...)
	if err != nil {
		return nil, errors.Wrap(err, "error creating Quantile")
	}

	return &IQR{quantile: quantile}, nil
}

// String returns a string representation of the metric.
func (i *IQR) String() string {
	name := "quantile.IQR"
	quantile := fmt.Sprintf("quantile:%v", i.quantile.String())
	return fmt.Sprintf("%s_{%s}", name, quantile)
}

// Push adds a number for calculating the interquartile range.
func (i *IQR) Push(x float64) error {
	err := i.quantile.Push(x)
	if err != nil {
		return errors.Wrapf(err, "error pushing %f to Quantile", x)
	}
	return nil
}

// Value returns the value of the interquartile range.
func (i *IQR) Value() (float64, error) {
	i.quantile.RLock()
	defer i.quantile.RUnlock()

	q25, err := i.quantile.Value(0.25)
	if err != nil {
		return 0, errors.Wrap(err, "error retrieving 1st quartile")
	}

	q75, err := i.quantile.Value(0.75)
	if err != nil {
		return 0, errors.Wrap(err, "error retrieving 3rd quartile")
	}

	return q75 - q25, nil
}

// Clear resets the metric.
func (i *IQR) Clear() {
	i.quantile.Clear()
}
