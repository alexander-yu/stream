# Stream

Stream is a Go library for online statistical algorithms. Provided statistics can be computed globally over an entire stream, or over a rolling window.

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
