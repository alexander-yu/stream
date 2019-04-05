package order

// Node is an interface that acts as a container for a value.
type Node interface {
	Value() float64
}

// Statistic is the interface required for any data structure that
// can provide order statistics.
type Statistic interface {
	Add(float64)
	Remove(float64)
	Size() int
	Select(int) Node
	Rank(float64) int
	Clear()
}

// Option is an optional argument which sets an optional field for creating an order.Statistic
type Option func(Statistic) error
