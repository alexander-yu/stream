package moment

import (
	"math"
	"sync"

	"github.com/pkg/errors"
	"github.com/workiva/go-datastructures/queue"
)

// Core is a struct that stores fundamental information for moments of a stream.
type Core struct {
	mux    sync.RWMutex
	mean   float64
	sums   []float64
	count  int
	window uint64
	queue  *queue.RingBuffer
}

// SetupMetric sets a Metric up with a core for consuming.
func SetupMetric(metric Metric) error {
	config := metric.Config()
	core, err := NewCore(config)
	if err != nil {
		return errors.Wrap(err, "error creating Core")
	}

	metric.Subscribe(core)
	return nil
}

// NewCore instantiates a Core struct based on a provided config.
func NewCore(config *CoreConfig) (*Core, error) {
	// validate config
	err := validateConfig(config)
	if err != nil {
		return nil, errors.Wrap(err, "error validating config")
	}

	// set defaults for any remaining unset fields
	config = setConfigDefaults(config)

	// initialize and create core
	c := &Core{}
	c.window = uint64(*config.Window)

	maxSum := -1
	for k := range config.Sums {
		if k > maxSum {
			maxSum = k
		}
	}
	c.sums = make([]float64, maxSum+1)

	c.queue = queue.NewRingBuffer(c.window)

	return c, nil
}

// Push adds a new value for a Core object to consume.
func (c *Core) Push(x float64) error {
	c.mux.Lock()
	defer c.mux.Unlock()
	return c.UnsafePush(x)
}

// UnsafePush adds a new value for a Core object to consume,
// but does not lock. This should only be used if the user
// plans to make use of the Lock()/Unlock() Core methods.
func (c *Core) UnsafePush(x float64) error {
	if c.window != 0 {
		if c.queue.Len() == c.window {
			tail, err := c.queue.Get()
			if err != nil {
				return errors.Wrap(err, "error popping item from queue")
			}

			c.remove(tail.(float64))
		}

		err := c.queue.Put(x)
		if err != nil {
			return errors.Wrapf(err, "error pushing %f to queue", x)
		}
	}

	c.add(x)
	return nil
}

// add updates the mean, count, and centralized power sums in an efficient
// and stable (numerically speaking) way, which allows for more accurate reporting
// of moments. See the following paper for details on the algorithm used:
// P. Pebay, T. B. Terriberry, H. Kolla, J. Bennett, Numerically stable, scalable
// formulas for parallel and online computation of higher-order multivariate central
// moments with arbitrary weights, Computational Statistics 31 (2016) 1305–1325.
func (c *Core) add(x float64) {
	c.count++
	count := float64(c.count)
	delta := x - c.mean
	c.mean += delta / count
	for k := len(c.sums) - 1; k >= 2; k-- {
		switch k {
		case 2:
			c.sums[k] += (count - 1) / count * math.Pow(delta, float64(k))
		default:
			c.sums[k] += (count - 1) / count * (math.Pow(count-1, float64(k-1)) + float64(sign(k))) / math.Pow(count, float64(k-1)) * math.Pow(delta, float64(k))
			for i := 1; i <= k-2; i++ {
				c.sums[k] += float64(binom(k, i)*sign(i)) * math.Pow(delta/count, float64(i)) * c.sums[k-i]
			}
		}
	}
}

// remove simply undoes the result of an add() call, and clears out the stats
// if we remove the last item of a window (only needed in the case where the
// window size is 1).
func (c *Core) remove(x float64) {
	c.count--
	if c.count > 0 {
		count := float64(c.count)
		c.mean -= (x - c.mean) / count
		delta := x - c.mean
		for k := range c.sums {
			switch k {
			case 0:
			case 1:
				continue
			case 2:
				c.sums[k] -= count / (count + 1) * math.Pow(delta, float64(k))
			default:
				c.sums[k] -= count / (count + 1) * (math.Pow(count, float64(k-1)) + float64(sign(k))) / math.Pow(count+1, float64(k-1)) * math.Pow(delta, float64(k))
				for i := 1; i <= k-2; i++ {
					c.sums[k] -= float64(binom(k, i)*sign(i)) * math.Pow(delta/(count+1), float64(i)) * c.sums[k-i]
				}
			}
		}
	} else {
		c.mean = 0
		for k := range c.sums {
			c.sums[k] = 0
		}
	}
}

// Count returns the number of values seen seen globally.
func (c *Core) Count() int {
	c.mux.RLock()
	defer c.mux.RUnlock()
	return c.count
}

// Mean returns the mean of values seem.
func (c *Core) Mean() (float64, error) {
	c.mux.RLock()
	defer c.mux.RUnlock()

	if c.count == 0 {
		return 0, errors.New("no values seen yet")
	}

	return c.mean, nil
}

// Sum returns the kth-power centralized sum of values seen.
// In other words, this returns the kth power sum of the differences
// of the values seen from their mean.
func (c *Core) Sum(k int) (float64, error) {
	c.mux.RLock()
	defer c.mux.RUnlock()

	if c.count == 0 {
		return 0, errors.New("no values seen yet")
	}

	if k <= 0 || k >= len(c.sums) {
		return 0, errors.Errorf("%d is not a tracked power sum", k)
	}

	return c.sums[k], nil
}

// Clear clears all stats being tracked.
func (c *Core) Clear() {
	c.mux.Lock()
	defer c.mux.Unlock()

	for k := range c.sums {
		c.sums[k] = 0
	}

	c.count = 0
	c.queue.Dispose()
	c.queue = queue.NewRingBuffer(c.window)
}

// RLock locks the core internals for reading.
func (c *Core) RLock() {
	c.mux.RLock()
}

// RUnlock undoes a single RLock call.
func (c *Core) RUnlock() {
	c.mux.RUnlock()
}

// Lock locks the core internals for writing.
func (c *Core) Lock() {
	c.mux.Lock()
}

// Unlock undoes a Lock call.
func (c *Core) Unlock() {
	c.mux.Unlock()
}
