package stream

import (
	"sync"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
)

// SimpleAggregateMetric is a wrapper metric that tracks multiple single-value metrics simultaneously.
// Note that it simply stores multiple metrics and pushes to all of them; this can be inefficient
// for metrics that could make use of shared data.
type SimpleAggregateMetric struct {
	metrics []Metric
	mux     sync.Mutex
}

// NewSimpleAggregateMetric instantiates an SimpleAggregateMetric struct.
func NewSimpleAggregateMetric(metrics ...Metric) *SimpleAggregateMetric {
	return &SimpleAggregateMetric{metrics: metrics}
}

// Push adds a new value for the metrics to consume.
func (s *SimpleAggregateMetric) Push(x float64) error {
	s.mux.Lock()
	defer s.mux.Unlock()
	var errs []error
	var mux sync.Mutex
	var wg sync.WaitGroup

	for _, metric := range s.metrics {
		wg.Add(1)
		go func(metric Metric) {
			defer wg.Done()
			err := metric.Push(x)
			if err != nil {
				mux.Lock()
				errs = append(errs, err)
				mux.Unlock()
			}
		}(metric)
	}

	wg.Wait()

	if len(errs) != 0 {
		var result *multierror.Error
		for _, err := range errs {
			result = multierror.Append(result, err)
		}
		return errors.Wrapf(result, "error pushing %f to metrics", x)
	}

	return nil
}

// Values returns the values of the metrics; in particular, it returns
// a map of strings to values, where the strings are the string
// representations of each metric (i.e. the result of calling String()).
func (s *SimpleAggregateMetric) Values() (map[string]float64, error) {
	s.mux.Lock()
	defer s.mux.Unlock()
	values := map[string]float64{}
	var errs []error
	var mux sync.Mutex
	var wg sync.WaitGroup

	for _, metric := range s.metrics {
		wg.Add(1)
		go func(metric Metric) {
			defer wg.Done()
			val, err := metric.Value()

			mux.Lock()
			if err != nil {
				errs = append(errs, err)
			} else {
				values[metric.String()] = val
			}
			mux.Unlock()
		}(metric)
	}

	wg.Wait()

	if len(errs) != 0 {
		var result *multierror.Error
		for _, err := range errs {
			result = multierror.Append(result, err)
		}
		return nil, errors.Wrap(result, "error retrieving values from metrics")
	}

	return values, nil
}
