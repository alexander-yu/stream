# Stream

[![GoDoc](http://godoc.org/github.com/alexander-yu/stream?status.svg)](http://godoc.org/github.com/alexander-yu/stream)
[![Build Status](https://travis-ci.org/alexander-yu/stream.svg?branch=master)](https://travis-ci.org/alexander-yu/stream)
[![Go Report Card](https://goreportcard.com/badge/github.com/alexander-yu/stream)](https://goreportcard.com/report/github.com/alexander-yu/stream)
[![codecov](https://codecov.io/gh/alexander-yu/stream/branch/master/graph/badge.svg)](https://codecov.io/gh/alexander-yu/stream)
[![GitHub license](https://img.shields.io/github/license/alexander-yu/stream.svg)](https://github.com/alexander-yu/stream/blob/master/LICENSE)

Stream is a Go library for online statistical algorithms. Provided statistics can be computed globally over an entire stream, or over a rolling window.

## Table of Contents

- [Stream](#stream)
  - [Table of Contents](#table-of-contents)
  - [Installation](#installation)
  - [Example Usage](#example-usage)
  - [Statistics](#statistics)
    - [Quantile](#quantile)
      - [Quantile](#quantile-1)
      - [Median](#median)
      - [IQR](#iqr)
      - [HeapMedian](#heapmedian)
    - [Min/Max](#minmax)
      - [Min](#min)
      - [Max](#max)
    - [Moment-Based Statistics](#moment-based-statistics)
      - [Mean](#mean)
      - [EWMA](#ewma)
      - [Moment](#moment)
      - [EWMMoment](#ewmmoment)
      - [Std](#std)
      - [Skewness](#skewness)
      - [Kurtosis](#kurtosis)
      - [Core (Univariate)](#core-univariate)
    - [Joint Distribution Statistics](#joint-distribution-statistics)
      - [Covariance](#covariance)
      - [Correlation](#correlation)
      - [Autocorrelation](#autocorrelation)
      - [Core (Multivariate)](#core-multivariate)
    - [Aggregate Statistics](#aggregate-statistics)
      - [SimpleAggregateMetric](#simpleaggregatemetric)
      - [SimpleJointAggregateMetric](#simplejointaggregatemetric)

## Installation

Use `go get`:

```bash
go get github.com/alexander-yu/stream
```

## Example Usage

In-depth examples are provided in the [examples](https://github.com/alexander-yu/stream/tree/master/examples) directory, but a small taste is provided below:

```go
// tracks the autocorrelation over a
// rolling window of size 15 and lag of 5
autocorr, err := joint.NewAutocorrelation(5, 15)
// handle err

// all metrics in the joint package must be passed
// through joint.Init in order to consume values
err = joint.Init(autocorr)
// handle err

// tracks the global median using a pair of heaps
median, err := quantile.NewHeapMedian(0)
// handle err

for i := 0., i < 100; i++ {
    err = autocorr.Push(i)
    // handle err

    err = median.Push(i)
    // handle err
}

autocorrVal, err := autocorr.Value()
// handle err

medianVal, err := median.Value()
// handle err

fmt.Println("%s: %f", autocorr.String(), autocorrVal)
fmt.Println("%s: %f", median.String(), medianVal)
```

## Statistics

For time/space complexity details on the algorithms listed below, see [here](complexity.md).

### [Quantile](https://godoc.org/github.com/alexander-yu/stream/quantile)

#### Quantile

Quantile keeps track of a given quantile of a stream with an [order statistic tree](https://en.wikipedia.org/wiki/Order_statistic_tree) as the underlying data structure. Quantile can calculate the global quantile of a stream, or over a rolling window. You can also configure which implementation to use for the underlying order statistic tree (see the [godoc](https://godoc.org/github.com/alexander-yu/stream/quantile#NewQuantile) entry for details), as well as which interpolation method to use in the case that the quantile actually lies in between two elements. For now only [AVL trees](https://en.wikipedia.org/wiki/AVL_tree) and [red black trees](https://en.wikipedia.org/wiki/Red-black_tree) are supported.

#### Median

Median keeps track of the median of a stream; this is simply a convenient wrapper over [Quantile](#Quantile), that automatically sets the quantile to be 0.5 and the interpolation method to be the midpoint method.

#### IQR

IQR keeps track of the [interquartile range](https://en.wikipedia.org/wiki/Interquartile_range) of a stream; this is simply a convenient wrapper over [Quantile](#Quantile), that retrieves the 1st and 3rd quartiles and sets the interpolation method to be the midpoint method.

#### HeapMedian

HeapMedian keeps track of the median of a stream with a pair of [heaps](https://en.wikipedia.org/wiki/Heap_(data_structure)). In particular, it uses a max-heap and a min-heap to keep track of elements below and above the median, respectively. HeapMedian can calculate the global median of a stream, or over a rolling window.

### [Min/Max](https://godoc.org/github.com/alexander-yu/stream/minmax)

#### Min

Min keeps track of the minimum of a stream; it can track either the global minimum, or over a rolling window.

#### Max

Max keeps track of the maximum of a stream; it can track either the global maximum, or over a rolling window.

### [Moment-Based Statistics](https://godoc.org/github.com/alexander-yu/stream/moment)

#### Mean

Mean keeps track of the mean of a stream; it can track either the global mean, or over a rolling window.

#### EWMA

EWMA keeps track of the global [exponentially weighted moving average](https://en.wikipedia.org/wiki/Moving_average#Exponential_moving_average).

#### Moment

Moment keeps track of the `k`-th sample [central moment](https://en.wikipedia.org/wiki/Central_moment); it can track either the global moment, or over a rolling window.

#### EWMMoment

EWMMoment keeps track of the global exponentially weighted moving central moment. This uses the exponentially weighted moving average as its center of mass, and uses the same exponential weights for its power terms.

#### Std

Std keeps track of the sample [standard deviation](https://en.wikipedia.org/wiki/Standard_deviation) of a stream; it can track either the global standard deviation, or over a rolling window. To track the sample [variance](https://en.wikipedia.org/wiki/Variance) instead, you should use [Moment](#Moment), i.e.

```go
variance := NewMoment(2, window)
```

#### Skewness

Skewness keeps track of the sample [skewness](https://en.wikipedia.org/wiki/Skewness) of a stream (in particular, the [adjusted Fisher-Pearson standardized moment coefficient](https://en.wikipedia.org/wiki/Skewness#Sample_skewness)); it can track either the global skewness, or over a rolling window.

#### Kurtosis

Kurtosis keeps track of the sample [kurtosis](https://en.wikipedia.org/wiki/Kurtosis) of a stream (in particular, the [sample excess kurtosis](https://en.wikipedia.org/wiki/Kurtosis#Sample_kurtosis)); it can track either the global kurtosis, or over a rolling window.

#### Core (Univariate)

Core is the struct powering all of the statistics in the `stream/moment` subpackage; it keeps track of a pre-configured set of centralized `k`-th power sums of a stream in an efficient, numerically stable way; it can track either the global sums, or over a rolling window.

To configure which sums to track, you'll need to instantiate a `CoreConfig` struct and provide it to `NewCore`:

```go
config := &moment.CoreConfig{
    Sums: SumsConfig{
        2: true, // tracks the sum of squared differences
        3: true, // tracks the sum of cubed differences
    },
    Window: stream.IntPtr(0), // track global sums
}
core, err := NewCore(config)
```

See the [godoc](https://godoc.org/github.com/alexander-yu/stream/moment#Core) entry for more details on Core's methods.

### [Joint Distribution Statistics](https://godoc.org/github.com/alexander-yu/stream/joint)

#### Covariance

Covariance keeps track of the sample [covariance](https://en.wikipedia.org/wiki/Covariance) of a stream; it can track either the global covariance, or over a rolling window.

#### Correlation

Correlation keeps track of the sample [correlation](https://en.wikipedia.org/wiki/Correlation) of a stream (in particular, the [sample Pearson correlation coefficient](https://en.wikipedia.org/wiki/Pearson_correlation_coefficient#For_a_sample)); it can track either the global correlation, or over a rolling window.

#### Autocorrelation

Autocorrelation keeps track of the sample [autocorrelation](https://en.wikipedia.org/wiki/Autocorrelation) of a stream (in particular, the [sample autocorrelation](https://en.wikipedia.org/wiki/Autocorrelation#Estimation)) for a given lag; it can track either the global autocorrelation, or over a rolling window.

#### Core (Multivariate)

Core is the struct powering all of the statistics in the `stream/joint` subpackage; it keeps track of a pre-configured set of joint centralized power sums of a stream in an efficient, numerically stable way; it can track either the global sums, or over a rolling window.

To configure which sums to track, you'll need to instantiate a `CoreConfig` struct and provide it to `NewCore`:

```go
config := &joint.CoreConfig{
    Sums: SumsConfig{
        {1, 1}, // tracks the joint sum of differences
        {2, 0}, // tracks the sum of squared differences of variable 1
    },
    Vars: stream.IntPtr(2), // declares that there are 2 variables to track (optional if Sums is set)
    Window: stream.IntPtr(0), // track global sums
}
core, err := NewCore(config)
```

See the [godoc](https://godoc.org/github.com/alexander-yu/stream/joint#Core) entry for more details on Core's methods.

### [Aggregate Statistics](https://godoc.org/github.com/alexander-yu/stream/aggregate)

#### SimpleAggregateMetric

SimpleAggregateMetric is a convenience wrapper that stores multiple univariate metrics and will push a value to all metrics simultaneously; instead of returning a single scalar, it returns a map of metrics to their corresponding values.

#### SimpleJointAggregateMetric

SimpleJointAggregateMetric is a convenience wrapper that stores multiple multivariate metrics and will push a value to all metrics simultaneously; instead of returning a single scalar, it returns a map of metrics to their corresponding values.
