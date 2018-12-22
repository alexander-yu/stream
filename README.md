# Stream

[![GoDoc](http://godoc.org/github.com/alexander-yu/stream?status.svg)](http://godoc.org/github.com/alexander-yu/stream)
[![Build Status](https://travis-ci.org/alexander-yu/stream.svg?branch=master)](https://travis-ci.org/alexander-yu/stream)
[![codecov](https://codecov.io/gh/alexander-yu/stream/branch/master/graph/badge.svg)](https://codecov.io/gh/alexander-yu/stream)
[![Go Report Card](https://goreportcard.com/badge/github.com/alexander-yu/stream)](https://goreportcard.com/report/github.com/alexander-yu/stream)

Stream is a Go library for online statistical algorithms. Provided statistics can be computed globally over an entire stream, or over a rolling window.

## Table of Contents
- [Installation](#Installation)
- [Introduction](#Introduction)
- [Example](#Example)
- [Statistics](#Statistics)
    - [Median](#Median)
        - [AVLMedian](#AVLMedian)
        - [HeapMedian](#HeapMedian)
    - [Min/Max](#Min/Max)
        - [Min](#Min)
        - [Max](#Max)
    - [Moment-Based Statistics](#Moment-Based-Statistics)
        - [Mean](#Mean)
        - [Moment](#Moment)
        - [Variance](#Variance)
        - [Std](#Std)
        - [Skewness](#Skewness)
        - [Kurtosis](#Kurtosis)
        - [Core](#Core)

## Installation
Use `go get`:

```
$ go get github.com/alexander-yu/stream
```

## Introduction
Every metric satisfies the following interface:
```go
type Metric interface {
    Push(float64) error
    Value() (float64, error)
}
```
The `Push` method consumes a numeric value (i.e. `float64`), and returns an error if one was encountered while pushing. The `Value` method returns the value of the metric at that given point in time, or an error if one was encountered when attempting to retrieve the value.

## Example
```go
package main

import (
    "fmt"

    "github.com/alexander-yu/stream/moment"
)

func main() {
    mean, err := moment.NewMean(5)
    if err != nil {
        fmt.Printf("error getting mean: %+v", err)
        return
    }

    for i := 0.; i < 100; i++ {
        err := mean.Push(i)
        if err != nil {
            fmt.Printf("error pushing %f: %+v", i, err)
            return
        }
    }

    val, err := mean.Value()
    if err != nil {
        fmt.Printf("error getting value: %+v", err)
        return
    }

    fmt.Println("mean:", val)
}
```

## Statistics

### [Median](https://godoc.org/github.com/alexander-yu/stream/median)

#### AVLMedian
AVLMedian keeps track of the median of a stream with an [AVL tree](https://en.wikipedia.org/wiki/AVL_tree) as the underlying data structure; more specifically, it uses an AVL tree that is also an [order statistic tree](https://en.wikipedia.org/wiki/Order_statistic_tree). AVLMedian can calculate the global median of a stream, or over a rolling window.

Let `n` be the size of the window, or the stream if tracking the global median. Then we have the following complexities:

| Push (time) | Value (time) | Space  |
| :---------: | :----------: | :----: |
| `O(log n)`  | `O(log n)`   | `O(n)` |

#### HeapMedian
HeapMedian keeps track of the global median of a stream with a pair of [heaps](https://en.wikipedia.org/wiki/Heap_(data_structure)). In particular, it uses a max-heap and a min-heap to keep track of elements below and above the median, respectively.

**Note:** HeapMedian does not support calculating medians over a rolling window, due to the non-constant time complexity required to remove expired values from the heaps as they leave the window.

Let `n` be the size of the stream. Then we have the following complexities:

| Push (time) | Value (time) | Space  |
| :---------: | :----------: | :----: |
| `O(log n)`  | `O(log n)`   | `O(n)` |


### [Min/Max](https://godoc.org/github.com/alexander-yu/stream/minmax)

#### Min
Min keeps track of the minimum of a stream; it can track either the global minimum, or over a rolling window.

Let `n` be the size of the window, or the stream if tracking the global minimum. Then we have the following complexities:

| Push (time)        | Value (time)       | Space                         |
| :----------------: | :----------------: | :---------------------------: |
| `O(1)` (amortized) | `O(1)` (amortized) | `O(1)` if global, else `O(n)` |

#### Max
Max keeps track of the maximum of a stream; it can track either the global maximum, or over a rolling window.

Let `n` be the size of the window, or the stream if tracking the global maximum. Then we have the following complexities:

| Push (time)        | Value (time)       | Space                         |
| :----------------: | :----------------: | :---------------------------: |
| `O(1)` (amortized) | `O(1)` (amortized) | `O(1)` if global, else `O(n)` |


### [Moment-Based Statistics](https://godoc.org/github.com/alexander-yu/stream/moment)

#### Mean
Mean keeps track of the mean of a stream; it can track either the global mean, or over a rolling window.

| Push (time) | Value (time) | Space  |
| :---------: | :----------: | :----: |
| `O(1)`      | `O(1)`       | `O(1)` |

#### Moment
Moment keeps track of the `k`-th sample [central moment](https://en.wikipedia.org/wiki/Central_moment); it can track either the global moment, or over a rolling window.

| Push (time) | Value (time) | Space  |
| :---------: | :----------: | :----: |
| `O(1)`      | `O(1)`       | `O(1)` |

#### Variance
Variance keeps track of the sample [variance](https://en.wikipedia.org/wiki/Variance) of a stream; it can track either the global variance, or over a rolling window.

| Push (time) | Value (time) | Space  |
| :---------: | :----------: | :----: |
| `O(1)`      | `O(1)`       | `O(1)` |

#### Std
Std keeps track of the sample [standard deviation](https://en.wikipedia.org/wiki/Standard_deviation) of a stream; it can track either the global standard deviation, or over a rolling window.

| Push (time) | Value (time) | Space  |
| :---------: | :----------: | :----: |
| `O(1)`      | `O(1)`       | `O(1)` |

#### Skewness
Skewness keeps track of the sample [skewness](https://en.wikipedia.org/wiki/Skewness) of a stream (in particular, the [adjusted Fisher-Pearson standardized moment coefficient](https://en.wikipedia.org/wiki/Skewness#Sample_skewness)); it can track either the global skewness, or over a rolling window.

| Push (time) | Value (time) | Space  |
| :---------: | :----------: | :----: |
| `O(1)`      | `O(1)`       | `O(1)` |

#### Kurtosis
Kurtosis keeps track of the sample [kurtosis](https://en.wikipedia.org/wiki/Kurtosis) of a stream (in particular, the [sample excess kurtosis](https://en.wikipedia.org/wiki/Kurtosis#Sample_kurtosis)); it can track either the global kurtosis, or over a rolling window.

| Push (time) | Value (time) | Space  |
| :---------: | :----------: | :----: |
| `O(1)`      | `O(1)`       | `O(1)` |

#### Core
Core is the struct powering all of the statistics in the `stream/moment` subpackage; it keeps track of a pre-configured set of `k`-th raw moments of a stream; it can track either the global moments, or over a rolling window. This can be used to implement any custom moment-based statistic you desire that isn't part of this package.

To configure which raw moments to track, you'll need to instantiate a `CoreConfig` struct and provide it to `NewCore`:

```go
config := &moment.CoreConfig{
    Sums: SumsConfig{
        1: true, // tracks the 1st raw moment, i.e. the sum
        2: true, // tracks the 2nd raw moment, i.e. the sum of squares
    },
    Window: stream.IntPtr(0), // track global moments
}
core, err := NewCore(config)
```

See the [godoc](https://godoc.org/github.com/alexander-yu/stream/moment#Core) entry for more details on Core's methods (note it does not satisfy the Metric interface, since it is capable of storing multiple values).

| Push (time) | Sum (time) | Count (time) | Space  |
| :---------: | :--------: | :----------: | :----: |
| `O(1)`      | `O(1)`     | `O(1)`       | `O(1)` |
