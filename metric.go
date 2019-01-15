package stream

// Every metric satisfies one of the following interfaces below.

// Metric is the interface for a metric that consumes from a stream.
// Metric is the standard interface for most metrics; in particular
// for those that consume single numeric values at a time.
type Metric interface {
	Push(float64) error
	Value() (float64, error)
	String() string
	Clear()
}

// AggregateMetric is the interface for a metric that tracks multiple single-value metrics simultaneously.
// Values() returns a map of metrics to their corresponding values at that given
// time. The keys are the string representations of the metrics (by calling the String() method).
type AggregateMetric interface {
	Push(float64) error
	Values() (map[string]interface{}, error)
}

// JointMetric is the interface for a metric that tracks joint statistics from a stream.
type JointMetric interface {
	Push(...float64) error
	Value() (float64, error)
}
