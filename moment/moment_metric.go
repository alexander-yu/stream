package moment

import "github.com/alexander-yu/stream"

// Metric is the interface for a metric tracks the moment of a stream.
type Metric interface {
	stream.Metric
	Subscribe(*Core)
	Config() *CoreConfig
}
