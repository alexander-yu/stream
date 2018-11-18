package sample

import (
	"math/rand"
)

// Reservoir is a struct for performing reservoir sampling on a stream.
type Reservoir struct {
	Size   int
	sample []interface{}
	count  int
}

// Push consumes a value to perform reservoir sampling.
func (r *Reservoir) Push(x interface{}) {
	r.count++
	if r.count <= r.Size {
		r.sample = append(r.sample, x)
	} else if index := rand.Intn(r.count); index < r.Size {
		r.sample[index] = x
	}
}

// Sample returns a copied slice of the obtained sample.
func (r *Reservoir) Sample() []interface{} {
	sample := make([]interface{}, len(r.sample))
	copy(sample, r.sample)
	return sample
}
