package joint

import "github.com/alexander-yu/stream"

// Metric is the interface for a metric that tracks joint statistics of a stream.
// Any Metric that will actually be consuming values (i.e. will have its Push
// method called) needs to be passed into the Init() method, which sets it up
// with a Core for consuming values and keeping track of centralized sums.
type Metric interface {
	stream.SimpleJointMetric
	SetCore(*Core)
	Config() *CoreConfig
}
