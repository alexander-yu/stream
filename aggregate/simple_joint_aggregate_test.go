package aggregate

import (
	"fmt"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/alexander-yu/stream"
)

type mockJointMetric struct {
	vals    [][]float64
	val     float64
	pushErr bool
	valErr  bool
}

func (m *mockJointMetric) String() string {
	return fmt.Sprintf("mockJointMetric_val:%f", m.val)
}

func (m *mockJointMetric) Push(xs ...float64) error {
	if m.pushErr {
		return errors.Errorf("error pushing %v", xs)
	}

	m.vals = append(m.vals, xs)
	return nil
}

func (m *mockJointMetric) Value() (float64, error) {
	if m.valErr {
		return 0, errors.New("error retrieving value")
	}

	return m.val, nil
}

func (m *mockJointMetric) Clear() {
	m.vals = [][]float64{}
}

func TestNewSimpleJointAggregateMetric(t *testing.T) {
	metric1 := &mockJointMetric{}
	metric2 := &mockJointMetric{}
	metric := NewSimpleJointAggregateMetric(metric1, metric2)
	assert.Equal(t, []stream.SimpleJointMetric{metric1, metric2}, metric.metrics)
}

func TestSimpleJointAggregateMetricPush(t *testing.T) {
	t.Run("pass: pushes value to each metric", func(t *testing.T) {
		metric1 := &mockJointMetric{}
		metric2 := &mockJointMetric{}
		metric := NewSimpleJointAggregateMetric(metric1, metric2)

		for i := 0.; i < 5; i++ {
			err := metric.Push(i, i*i)
			require.NoError(t, err)
		}

		expected := [][]float64{{0, 0}, {1, 1}, {2, 4}, {3, 9}, {4, 16}}
		assert.Equal(t, expected, metric1.vals)
		assert.Equal(t, expected, metric2.vals)
	})

	t.Run("fail: returns error if any Push() call fails", func(t *testing.T) {
		metric1 := &mockJointMetric{}
		metric2 := &mockJointMetric{pushErr: true}
		metric3 := &mockJointMetric{pushErr: true}
		metric := NewSimpleJointAggregateMetric(metric1, metric2, metric3)

		err := metric.Push(0., 0.)
		assert.EqualError(t, err, fmt.Sprintf(
			"error pushing %v to metrics: 2 errors occurred:\n\t* error pushing %v\n\t* error pushing %v\n\n",
			[]float64{0, 0},
			[]float64{0, 0},
			[]float64{0, 0},
		))
	})
}

func TestSimpleJointAggregateMetricValue(t *testing.T) {
	t.Run("pass: retrieves value of each metric", func(t *testing.T) {
		metric1 := &mockJointMetric{val: 1}
		metric2 := &mockJointMetric{val: 2}
		metric := NewSimpleJointAggregateMetric(metric1, metric2)

		values, err := metric.Values()
		require.NoError(t, err)

		expected := map[string]float64{
			metric1.String(): metric1.val,
			metric2.String(): metric2.val,
		}

		assert.Equal(t, expected, values)
	})

	t.Run("fail: returns error if any Value() call fails", func(t *testing.T) {
		metric1 := &mockJointMetric{val: 1}
		metric2 := &mockJointMetric{val: 2, valErr: true}
		metric3 := &mockJointMetric{val: 3, valErr: true}
		metric := NewSimpleJointAggregateMetric(metric1, metric2, metric3)

		_, err := metric.Values()
		assert.EqualError(t, err, "error retrieving values from metrics: 2 errors occurred:\n\t* error retrieving value\n\t* error retrieving value\n\n")

	})
}

func TestSimpleJointAggregateMetricClear(t *testing.T) {
	metric1 := &mockJointMetric{}
	metric2 := &mockJointMetric{}
	metric := NewSimpleJointAggregateMetric(metric1, metric2)

	for i := 0.; i < 5; i++ {
		err := metric.Push(i)
		require.NoError(t, err)
	}

	metric.Clear()
	assert.Equal(t, 0, len(metric1.vals))
	assert.Equal(t, 0, len(metric2.vals))
}
