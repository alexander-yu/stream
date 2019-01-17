package stream

// Every metric satisfies one of the following interfaces below.

// Metric is the interface for a metric that consumes from a stream.
// Metric is the standard interface for most metrics; in particular
// for those that consume single numeric values at a time. There is no
// Value method for this interface, allowing implementations to roll
// custom value methods.
type Metric interface {
	Push(float64) error
	String() string
	Clear()
}

// SimpleMetric is the interface for a Metric that returns a singular value.
type SimpleMetric interface {
	Metric
	Value() (float64, error)
}

// AggregateMetric is the interface for a metric that tracks multiple single-value metrics simultaneously.
// Values() returns a map of metrics to their corresponding values at that given
// time. The keys are the string representations of the metrics (by calling the String() method).
type AggregateMetric interface {
	Push(float64) error
	Values() (map[string]interface{}, error)
	Clear()
}

// JointMetric is the interface for a metric that tracks joint statistics from a stream.
// There is no Value method for this interface, allowing implementations to roll
// custom value methods.
type JointMetric interface {
	Push(...float64) error
	Clear()
}

// SimpleJointMetric is the interface for a JointMetric that returns a singular value.
type SimpleJointMetric interface {
	JointMetric
	Value() (float64, error)
}
