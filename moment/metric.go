package moment

import "github.com/alexander-yu/stream"

// Metric is the interface for a metric that tracks the moment of a stream.
type Metric interface {
	stream.SimpleMetric
	Subscribe(*Core)
	Config() *CoreConfig
}
