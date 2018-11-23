package stream

import (
	"math"

	"github.com/pkg/errors"
	"github.com/workiva/go-datastructures/queue"
)

// Core is a struct that stores fundamental information for stats collection on a stream.
type Core struct {
	sums        map[int]float64
	count       int
	min         float64
	max         float64
	window      uint64
	queue       *queue.RingBuffer
	pushMetrics []Metric
}

// NewCore creates a new Core struct based on a provided config as well as any configs
// for optionally provided metrics and subscribes those metrics to the Core.
func NewCore(config *CoreConfig, metrics ...Metric) (*Core, error) {
	// validate configs
	configs := []*CoreConfig{config}
	for _, metric := range metrics {
		config := metric.Config()
		err := validateConfig(config)
		if err != nil {
			return nil, errors.Wrap(err, "error validating config")
		}
		configs = append(configs, config)
	}

	// merge metric configs and set defaults for any remaining unset fields
	config, err := MergeConfigs(configs...)
	if err != nil {
		return nil, errors.Wrap(err, "error merging metric configs")
	}

	config = setConfigDefaults(config)

	// initialize and create Core
	c := &Core{min: math.Inf(1), max: math.Inf(-1)}
	c.window = uint64(*config.Window)
	c.sums = map[int]float64{}
	for k := range config.Sums {
		c.sums[k] = 0
	}
	c.queue = queue.NewRingBuffer(c.window)
	c.pushMetrics = config.PushMetrics

	for _, metric := range metrics {
		metric.Subscribe(c)
	}

	return c, nil
}

// Push adds a new value for a Core object to consume.
func (c *Core) Push(x float64) error {
	// Push value to all push metrics after completion
	defer func() {
		for _, metric := range c.pushMetrics {
			metric.Push(x)
		}
	}()

	if c.window != 0 {
		err := c.queue.Put(x)
		if err != nil {
			return errors.Wrap(err, "error pushing to metric")
		}

		if c.queue.Len() > c.window {
			tail, err := c.queue.Get()
			if err != nil {
				return errors.Wrap(err, "error popping from window queue")
			}

			c.count--

			tailVal := tail.(float64)
			for k := range c.sums {
				c.sums[k] -= math.Pow(tailVal, float64(k))
			}
		}
	}

	for k := range c.sums {
		c.sums[k] += math.Pow(x, float64(k))
	}

	c.count++
	c.min = math.Min(c.min, x)
	c.max = math.Max(c.max, x)
	return nil
}

// Count returns the number of values seen seen globally.
func (c *Core) Count() int {
	return c.count
}

// Min returns the min of values seen.
func (c *Core) Min() float64 {
	return c.min
}

// Max returns the max of values seen.
func (c *Core) Max() float64 {
	return c.max
}

// Sum returns the kth-power sum of values seen.
func (c *Core) Sum(k int) (float64, error) {
	if c.count == 0 {
		return 0, errors.New("no values seen yet")
	}

	if sum, ok := c.sums[k]; ok {
		return sum, nil
	}

	return 0, errors.Errorf("%d is not a tracked power sum", k)
}

// Clear clears all stats being tracked.
func (c *Core) Clear() {
	for k := range c.sums {
		c.sums[k] = 0
	}

	c.count = 0
	c.min = math.Inf(1)
	c.max = math.Inf(-1)
	c.queue.Dispose()
	c.queue = queue.NewRingBuffer(c.window)
}
