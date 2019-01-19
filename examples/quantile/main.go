package main

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/alexander-yu/stream"
	"github.com/alexander-yu/stream/quantile"
	"github.com/alexander-yu/stream/quantile/ost"
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
	ostMedian, err := quantile.NewOSTMedian(0, ost.RB)
	if err != nil {
		log.Fatal(err)
	}

	// tracks quantiles via an AVL tree over a rolling window of size 3
	// and with linear interpolation
	ostQuantile, err := quantile.NewOSTQuantile(
		&quantile.Config{
			Window:        stream.IntPtr(3),
			Interpolation: quantile.Linear.Ptr(),
		},
		ost.AVL,
	)
	if err != nil {
		log.Fatal(err)
	}

	// tracks the median over a rolling window of size 3 via a pair of heaps
	heapMedian, err := quantile.NewHeapMedian(3)
	if err != nil {
		log.Fatal(err)
	}

	metrics := []stream.Metric{ostMedian, ostQuantile, heapMedian}

	err = push(metrics)
	if err != nil {
		log.Fatal(err)
	}

	ostMedianVal, err := ostMedian.Value()
	if err != nil {
		log.Fatal(err)
	}

	// retrieve the 25% quantile
	ostQuantileVal, err := ostQuantile.Value(0.25)
	if err != nil {
		log.Fatal(err)
	}

	heapMedianVal, err := heapMedian.Value()
	if err != nil {
		log.Fatal(err)
	}

	values := map[string]float64{
		ostMedian.String():   ostMedianVal,
		ostQuantile.String(): ostQuantileVal,
		heapMedian.String():  heapMedianVal,
	}

	result, err := json.MarshalIndent(values, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(result))
}
