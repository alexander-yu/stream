package moment

import (
	"math"

	"github.com/pkg/errors"

	"github.com/alexander-yu/stream"
)

// Moment is a metric that tracks the kth sample central moment.
type Moment struct {
	k      int
	window int
	core   *stream.Core
}

// NewMoment instantiates a Moment struct that tracks the kth moment.
func NewMoment(k int, window int) (*Moment, error) {
	if k < 0 {
		return nil, errors.Errorf("%d is a negative moment", k)
	} else if window < 0 {
		return nil, errors.Errorf("%d is a negative window", window)
	}

	return &Moment{k: k, window: window}, nil
}

// Subscribe subscribes the Moment to a Core object.
func (m *Moment) Subscribe(c *stream.Core) {
	m.core = c
}

// Config returns the CoreConfig needed.
func (m *Moment) Config() *stream.CoreConfig {
	sums := stream.SumsConfig{}
	for i := 0; i <= m.k; i++ {
		sums[i] = true
	}

	return &stream.CoreConfig{
		Sums:   sums,
		Window: &m.window,
	}
}

// Push adds a new value for Moment to consume.
func (m *Moment) Push(x float64) error {
	err := m.core.Push(x)
	if err != nil {
		return errors.Wrap(err, "error pushing to core")
	}
	return nil
}

// Value returns the value of the kth sample central moment.
func (m *Moment) Value() (float64, error) {
	m.core.RLock()
	defer m.core.RUnlock()

	count := m.core.Count()
	if count == 0 {
		return 0, errors.New("no values seen yet")
	}

	if m.k == 0 {
		return 1., nil
	} else if m.k == 1 {
		return 0., nil
	}

	sum, err := m.core.Sum(1)
	if err != nil {
		return 0, errors.Wrap(err, "error retrieving 1-power sum")
	}
	mean := sum / float64(count)

	var moment float64
	for i := 0; i <= m.k; i++ {
		sum, err := m.core.Sum(i)
		if err != nil {
			return 0, errors.Wrapf(err, "error retrieving %d-power sum", i)
		}

		moment += float64(binom(m.k, i)*sign(m.k-i)) * math.Pow(mean, float64(m.k-i)) * sum
	}

	moment /= (float64(count) - 1.)

	return moment, nil
}
