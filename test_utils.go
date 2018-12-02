package stream

import (
	"fmt"
)

type mockMetric struct {
	id   int
	vals []float64
}

func (m *mockMetric) Subscribe(c *Core) {}

func (m *mockMetric) Config() *CoreConfig {
	return &CoreConfig{
		Sums: map[int]bool{
			-1: true,
			0:  true,
			1:  true,
			2:  true,
			3:  true,
			4:  true,
		},
		Window:      IntPtr(3),
		PushMetrics: []Metric{m},
	}
}

func (m *mockMetric) Push(x float64) error {
	m.vals = append(m.vals, x)
	return nil
}

func (m *mockMetric) Value() (float64, error) {
	return 0, nil
}

// TestData sets up a metric and populates a core with pushes for testing purposes.
func TestData(metric Metric) *Core {
	core, err := SetupMetric(metric)
	if err != nil {
		panic(fmt.Sprintf("%+v", err))
	}

	for i := 1.; i < 5; i++ {
		err := core.Push(i)
		if err != nil {
			panic(fmt.Sprintf("%+v", err))
		}
	}

	err = core.Push(8)
	if err != nil {
		panic(fmt.Sprintf("%+v", err))
	}

	fmt.Println(core.sums)

	return core
}
