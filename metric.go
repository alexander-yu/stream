package stream

// Metric is the interface for a metric that consumes from a stream.
type Metric interface {
	Push(float64) error
	Value() (float64, error)
}

// HashableMetric is the interface for a metric that returns a hash value as an ID.
type HashableMetric interface {
	Metric
	Hash() uint64
}

// AggregateMetric is the interface for a metric that tracks multiple metrics simultaneously.
type AggregateMetric interface {
	Push(float64) error
	Values() (map[uint64]float64, error)
}

// JointMetric is the interface for a metric that tracks joint statistics from a stream.
type JointMetric interface {
	Push(...float64) error
	Value() (float64, error)
}
