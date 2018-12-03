package sample

import (
	"math/rand"
	"sync"
)

// Reservoir is a struct for performing reservoir sampling on a stream.
type Reservoir struct {
	Size   int
	sample []interface{}
	count  int
	mux    sync.Mutex
}

// Push consumes a value to perform reservoir sampling.
func (r *Reservoir) Push(x interface{}) {
	r.mux.Lock()
	defer r.mux.Unlock()

	r.count++
	if r.count <= r.Size {
		r.sample = append(r.sample, x)
	} else if index := rand.Intn(r.count); index < r.Size {
		r.sample[index] = x
	}
}

// Sample returns a copied slice of the obtained sample.
func (r *Reservoir) Sample() []interface{} {
	r.mux.Lock()
	defer r.mux.Unlock()

	sample := make([]interface{}, len(r.sample))
	copy(sample, r.sample)
	return sample
}
