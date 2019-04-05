package main

import (
	"encoding/json"
	"fmt"

	multierror "github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/alexander-yu/stream"
	"github.com/alexander-yu/stream/quantile"
)

func push(metrics []stream.Metric) error {
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

func main() {
	// tracks the global median via a red-black tree
	median, err := quantile.NewMedian(0, quantile.ImplOption(quantile.RedBlack))
	if err != nil {
		log.Fatal(err)
	}

	// tracks quantiles via an AVL tree over a rolling window of size 3
	// and with linear interpolation
	avlQuantile, err := quantile.NewQuantile(
		3, quantile.InterpolationOption(quantile.Linear), quantile.ImplOption(quantile.AVL),
	)
	if err != nil {
		log.Fatal(err)
	}

	// tracks the median over a rolling window of size 3 via a pair of heaps
	heapMedian, err := quantile.NewHeapMedian(3)
	if err != nil {
		log.Fatal(err)
	}

	metrics := []stream.Metric{median, avlQuantile, heapMedian}

	err = push(metrics)
	if err != nil {
		log.Fatal(err)
	}

	medianVal, err := median.Value()
	if err != nil {
		log.Fatal(err)
	}

	// retrieve the 25% quantile
	avlQuantileval, err := avlQuantile.Value(0.25)
	if err != nil {
		log.Fatal(err)
	}

	heapMedianVal, err := heapMedian.Value()
	if err != nil {
		log.Fatal(err)
	}

	values := map[string]float64{
		median.String():      medianVal,
		avlQuantile.String(): avlQuantileval,
		heapMedian.String():  heapMedianVal,
	}

	result, err := json.MarshalIndent(values, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(result))
}
