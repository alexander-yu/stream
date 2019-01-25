package moment

import "github.com/alexander-yu/stream"

// Metric is the interface for a metric that tracks the moment of a stream.
// Any Metric that will actually be consuming values (i.e. will have its Push
// method called) needs to be passed into the Init() method, which sets it up
// with a Core for consuming values and keeping track of centralized sums.
type Metric interface {
	stream.SimpleMetric
	CoreWrapper
}

// CoreWrapper is the interface for an entity that wraps around a Core for stats.
// The methods below are required for setting up a Core for the wrapper.
type CoreWrapper interface {
	SetCore(*Core)
	Config() *CoreConfig
}
