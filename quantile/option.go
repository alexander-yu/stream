package quantile

import (
	"github.com/alexander-yu/stream/quantile/order"
	"github.com/pkg/errors"
)

// Option is an optional argument for creating quantile-based metrics,
// which sets an optional field for creating a Quantile.
type Option func(*Quantile) error

// ImplOption creates an option that sets the implementation for the
// underlying data structure.
func ImplOption(impl Impl, options ...order.Option) Option {
	return func(q *Quantile) error {
		if !impl.Valid() {
			return errors.Errorf("attempted to set invalid Impl %d", impl)
		}

		var err error
		q.statistic, err = impl.init(options...)
		return errors.Wrap(err, "error setting Impl")
	}
}

// InterpolationOption creates an option that sets the interpolation
// method.
func InterpolationOption(i Interpolation) Option {
	return func(q *Quantile) error {
		if !i.Valid() {
			return errors.Errorf("attempted to set invalid Interpolation %d", i)
		}

		q.interpolation = i
		return nil
	}
}
