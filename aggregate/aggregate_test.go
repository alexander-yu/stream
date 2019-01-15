package aggregate

import (
	"fmt"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/alexander-yu/stream"
)

type mockMetric struct {
	vals    []float64
	val     float64
	pushErr bool
	valErr  bool
}

func (m *mockMetric) String() string {
	return fmt.Sprintf("mockMetric_val:%f", m.val)
}

func (m *mockMetric) Push(x float64) error {
	if m.pushErr {
		return errors.Errorf("error pushing %f", x)
	}

	m.vals = append(m.vals, x)
	return nil
}

func (m *mockMetric) Value() (float64, error) {
	if m.valErr {
		return 0, errors.New("error retrieving value")
	}

	return m.val, nil
}

func TestNewSimpleAggregateMetric(t *testing.T) {
	metric1 := &mockMetric{}
	metric2 := &mockMetric{}
	metric := NewSimpleAggregateMetric(metric1, metric2)
	assert.Equal(t, []stream.Metric{metric1, metric2}, metric.metrics)
}

func TestSimpleAggregateMetricPush(t *testing.T) {
	t.Run("pass: pushes value to each metric", func(t *testing.T) {
		metric1 := &mockMetric{}
		metric2 := &mockMetric{}
		metric := NewSimpleAggregateMetric(metric1, metric2)

		for i := 0.; i < 5; i++ {
			err := metric.Push(i)
			require.NoError(t, err)
		}

		expected := []float64{0, 1, 2, 3, 4}
		assert.Equal(t, expected, metric1.vals)
		assert.Equal(t, expected, metric2.vals)
	})

	t.Run("fail: returns error if any Push() call fails", func(t *testing.T) {
		metric1 := &mockMetric{}
		metric2 := &mockMetric{pushErr: true}
		metric3 := &mockMetric{pushErr: true}
		metric := NewSimpleAggregateMetric(metric1, metric2, metric3)

		err := metric.Push(0.)
		assert.EqualError(t, err, fmt.Sprintf(
			"error pushing %f to metrics: 2 errors occurred:\n\t* error pushing %f\n\t* error pushing %f\n\n",
			0.,
			0.,
			0.,
		))
	})
}

func TestSimpleAggregateMetricValue(t *testing.T) {
	t.Run("pass: retrieves value of each metric", func(t *testing.T) {
		metric1 := &mockMetric{val: 1}
		metric2 := &mockMetric{val: 2}
		metric := NewSimpleAggregateMetric(metric1, metric2)

		values, err := metric.Values()
		require.NoError(t, err)

		expected := map[string]float64{
			metric1.String(): metric1.val,
			metric2.String(): metric2.val,
		}

		assert.Equal(t, expected, values)
	})

	t.Run("fail: returns error if any Value() call fails", func(t *testing.T) {
		metric1 := &mockMetric{val: 1}
		metric2 := &mockMetric{val: 2, valErr: true}
		metric3 := &mockMetric{val: 3, valErr: true}
		metric := NewSimpleAggregateMetric(metric1, metric2, metric3)

		_, err := metric.Values()
		assert.EqualError(t, err, "error retrieving values from metrics: 2 errors occurred:\n\t* error retrieving value\n\t* error retrieving value\n\n")

	})
}
