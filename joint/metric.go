package joint

import "github.com/alexander-yu/stream"

// Metric is the interface for a metric that tracks joint statistics of a stream.
type Metric interface {
	stream.Metric
	Subscribe(*Core)
	Config() *CoreConfig
}
