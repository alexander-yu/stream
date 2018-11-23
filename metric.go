package stream

// Metric is an interface for any metric that subscribes to a Core object.
type Metric interface {
	Subscribe(*Core)
	Config() *CoreConfig
	Push(float64)
	Value() (float64, error)
}
