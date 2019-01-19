package main

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/alexander-yu/stream/moment"
)

// all metrics in the moment package must be passed through moment.Init
// if you want to push values to them
func initialize(metrics []moment.Metric) error {
	var errs []error

	for _, metric := range metrics {
		err := moment.Init(metric)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) != 0 {
		var result *multierror.Error
		for _, err := range errs {
			result = multierror.Append(result, err)
		}
		return errors.Wrapf(result, "error initializing metrics")
	}

	return nil
}

func push(metrics []moment.Metric) error {
	var errs []error

	for _, metric := range metrics {
		for i := 0.; i < 100; i++ {
			err := metric.Push(i)
			if err != nil {
				errs = append(errs, err)
				break
			}
		}
	}

	if len(errs) != 0 {
		var result *multierror.Error
		for _, err := range errs {
			result = multierror.Append(result, err)
		}
		return errors.Wrapf(result, "error pushing to metrics")
	}

	return nil
}

func values(metrics []moment.Metric) (map[string]float64, error) {
	var errs []error

	result := map[string]float64{}
	for _, metric := range metrics {
		val, err := metric.Value()
		if err != nil {
			errs = append(errs, err)
			break
		}
		result[metric.String()] = val
	}

	if len(errs) != 0 {
		var result *multierror.Error
		for _, err := range errs {
			result = multierror.Append(result, err)
		}
		return nil, errors.Wrapf(result, "error retrieving values from metrics")
	}

	return result, nil
}

func main() {
	// tracks the global mean
	mean := moment.NewMean(0)
	// tracks the global standard deviation
	std := moment.NewStd(0)
	// tracks the variance over a rolling window of size 5
	variance := moment.NewMoment(2, 5)
	// tracks the skewness over a rolling window of size 5
	skewness := moment.NewSkewness(5)
	// tracks the kurtosis over a rolling window of size 5
	kurtosis := moment.NewKurtosis(5)

	metrics := []moment.Metric{mean, std, variance, skewness, kurtosis}

	err := initialize(metrics)
	if err != nil {
		log.Fatal(err)
	}

	err = push(metrics)
	if err != nil {
		log.Fatal(err)
	}

	values, err := values(metrics)
	if err != nil {
		log.Fatal(err)
	}

	result, err := json.MarshalIndent(values, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(result))
}
