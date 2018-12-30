package joint

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/alexander-yu/stream"
)

type mockMetric struct {
	vals [][]float64
	core *Core
}

func newMockMetric() *mockMetric {
	metric := &mockMetric{}
	err := SetupMetric(metric)
	if err != nil {
		panic(fmt.Sprintf("%+v", err))
	}

	return metric
}

func (m *mockMetric) Subscribe(c *Core) {
	m.core = c
}

func (m *mockMetric) Config() *CoreConfig {
	return &CoreConfig{
		Sums: SumsConfig{
			{2, 2},
		},
		Window: stream.IntPtr(3),
	}
}

func (m *mockMetric) Push(xs ...float64) error {
	m.vals = append(m.vals, xs)
	err := m.core.Push(xs...)
	if err != nil {
		return errors.Wrap(err, "error pushing to core")
	}

	return nil
}

func (m *mockMetric) Value() (float64, error) {
	return 0, nil
}

func testData(metric Metric) error {
	for i := 1.; i < 5; i++ {
		err := metric.Push(i, i*i)
		if err != nil {
			return errors.Wrapf(err, "failed to push (%f, %f) to metric", i, i*i)
		}
	}

	err := metric.Push(8., 64.)
	if err != nil {
		return errors.Wrapf(err, "failed to push (%f, %f) to metric", 8., 64.)
	}

	return nil
}