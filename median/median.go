package median

// Median is an interface that keeps track of a running median.
type Median interface {
	Push(x float64)
	Median() (float64, error)
}
