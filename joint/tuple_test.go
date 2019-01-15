package joint

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	testutil "github.com/alexander-yu/stream/util/test"
)

func TestAbs(t *testing.T) {
	m := Tuple{4, 3, 7, 0}
	assert.Equal(t, 14, m.abs())
}

func TestHash(t *testing.T) {
	m := Tuple{}
	assert.Equal(t, uint64(0), m.hash())

	m = Tuple{1, 4, 2, 3}
	assert.Equal(t, uint64(33700), m.hash())

	// different order should guarantee a different hash
	m = Tuple{1, 2, 3, 4}
	assert.Equal(t, uint64(31810), m.hash())

	// different object but same contents should guarantee the same hash
	m = Tuple{1, 2, 3, 4}
	assert.Equal(t, uint64(31810), m.hash())
}

func TestEq(t *testing.T) {
	m := Tuple{1, 2, 3}
	n := Tuple{1, 2, 3}
	assert.True(t, m.eq(n))

	n = Tuple{2, 1, 3}
	assert.False(t, m.eq(n))

	n = Tuple{1, 2, 3, 4}
	assert.False(t, m.eq(n))
}

func TestSub(t *testing.T) {
	t.Run("pass: returns difference of Tuples", func(t *testing.T) {
		m := Tuple{4, 3, 7, 5}
		n := Tuple{1, 3, 4, 0}
		diff, err := sub(m, n)
		require.NoError(t, err)
		assert.Equal(t, Tuple{3, 0, 3, 5}, diff)
	})

	t.Run("fail: returns error if Tuples have different lengths", func(t *testing.T) {
		m := Tuple{1, 2, 3}
		n := Tuple{1, 2, 3, 4}
		_, err := sub(m, n)
		assert.EqualError(t, err, fmt.Sprintf(
			"Tuples have different lengths: %d != %d",
			len(m),
			len(n),
		))

		m = Tuple{1, 2, 3, 4}
		n = Tuple{1, 2, 3}
		_, err = sub(m, n)
		assert.EqualError(t, err, fmt.Sprintf(
			"Tuples have different lengths: %d != %d",
			len(m),
			len(n),
		))
	})
}

func TestMultinom(t *testing.T) {
	t.Run("pass: returns multinomial coefficient", func(t *testing.T) {
		m := Tuple{4, 3, 7, 5}
		n := Tuple{1, 3, 4, 0}
		value, err := multinom(m, n)
		require.NoError(t, err)
		assert.Equal(t, 140, value)
	})

	t.Run("fail: returns error if Tuples have different lengths", func(t *testing.T) {
		m := Tuple{1, 2, 3}
		n := Tuple{1, 2, 3, 4}
		_, err := multinom(m, n)
		assert.EqualError(t, err, fmt.Sprintf(
			"Tuples have different lengths: %d != %d",
			len(m),
			len(n),
		))

		m = Tuple{1, 2, 3, 4}
		n = Tuple{1, 2, 3}
		_, err = multinom(m, n)
		assert.EqualError(t, err, fmt.Sprintf(
			"Tuples have different lengths: %d != %d",
			len(m),
			len(n),
		))
	})
}

func TestPow(t *testing.T) {
	t.Run("pass: returns multinomial expression", func(t *testing.T) {
		x := []float64{1., 2., 1.5, -1.}
		n := Tuple{3, 4, 2, 5}
		value, err := pow(x, n)
		require.NoError(t, err)
		testutil.Approx(t, -36., value)

		x = []float64{1., 0., 1.5, -1.}
		value, err = pow(x, n)
		require.NoError(t, err)
		testutil.Approx(t, 0., value)
	})

	t.Run("fail: returns error if Tuples have different lengths", func(t *testing.T) {
		x := []float64{1., 2., 1.5}
		n := Tuple{1, 2, 3, 4}
		_, err := pow(x, n)
		assert.EqualError(t, err, fmt.Sprintf(
			"Cannot exponentiate slice and Tuple with different lengths: %d != %d",
			len(x),
			len(n),
		))

		x = []float64{1., 2., 1.5, -1.}
		n = Tuple{1, 2, 3}
		_, err = pow(x, n)
		assert.EqualError(t, err, fmt.Sprintf(
			"Cannot exponentiate slice and Tuple with different lengths: %d != %d",
			len(x),
			len(n),
		))
	})
}

func TestIter(t *testing.T) {
	t.Run("pass: executes callback over all tuples in increasing order", func(t *testing.T) {
		tuple := Tuple{2, 3}
		expectedRuns := []uint64{}
		for j := 0; j <= tuple[1]; j++ {
			for i := 0; i <= tuple[0]; i++ {
				expectedRuns = append(expectedRuns, Tuple{i, j}.hash())
			}
		}

		runs := []uint64{}
		iter(tuple, false, func(xs ...int) {
			runs = append(runs, Tuple(xs).hash())
		})

		assert.Equal(t, expectedRuns, runs)
	})

	t.Run("pass: executes callback over all tuples in decreasing order", func(t *testing.T) {
		tuple := Tuple{2, 3}
		expectedRuns := []uint64{}
		for j := tuple[1]; j >= 0; j-- {
			for i := tuple[0]; i >= 0; i-- {
				expectedRuns = append(expectedRuns, Tuple{i, j}.hash())
			}
		}

		runs := []uint64{}
		iter(tuple, true, func(xs ...int) {
			runs = append(runs, Tuple(xs).hash())
		})

		assert.Equal(t, expectedRuns, runs)
	})
}
