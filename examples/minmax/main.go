package main

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/alexander-yu/stream"
	"github.com/alexander-yu/stream/minmax"
)

func push(metrics []stream.SimpleMetric) error {
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

func values(metrics []stream.SimpleMetric) (map[string]float64, error) {
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
	// tracks the min over a rolling window of size 5
	min, err := minmax.NewMin(5)
	if err != nil {
		log.Fatal(err)
	}

	// tracks the global max
	max, err := minmax.NewMax(0)
	if err != nil {
		log.Fatal(err)
	}

	metrics := []stream.SimpleMetric{min, max}

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
