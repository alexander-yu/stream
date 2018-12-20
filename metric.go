package stream

// Metric is the interface for a metric that consumes from a stream.
type Metric interface {
	Push(float64) error
	Value() (float64, error)
}
