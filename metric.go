package stream

import (
	"github.com/alexander-yu/stream/util/hashutil"
)

// Every metric satisfies one of the following interfaces below.

// Metric is the interface for a metric that consumes from a stream.
// Metric is the standard interface for most metrics; in particular
// for those that consume single numeric values at a time.
type Metric interface {
	Push(float64) error
	Value() (float64, error)
}

// MappableMetric is the interface for a metric that can be stored in
// a map. It is effectively a Metric that also implements the Mappable
// interface.
type MappableMetric interface {
	Metric
	hashutil.Mappable
}

// AggregateMetric is the interface for a metric that tracks multiple metrics simultaneously.
// Values() returns a map of metrics to their corresponding values at that given
// time. The sub-metrics tracked are MappableMetrics, in order to distinguish the
// different metric values.
type AggregateMetric interface {
	Push(float64) error
	Values() (hashutil.Map, error)
}

// JointMetric is the interface for a metric that tracks joint statistics from a stream.
type JointMetric interface {
	Push(...float64) error
	Value() (float64, error)
}
