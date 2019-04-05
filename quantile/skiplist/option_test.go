package skiplist

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/alexander-yu/stream/quantile/ost/rb"
	testutil "github.com/alexander-yu/stream/util/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMaxLevelOption(t *testing.T) {
	t.Run("fail: non-skiplist is invalid", func(t *testing.T) {
		err := MaxLevelOption(5)(&rb.Tree{})
		testutil.ContainsError(t, err, "attempted to set max level on a non-skiplist")
	})

	t.Run("fail: max level < 1 is invalid", func(t *testing.T) {
		skiplist, err := New()
		require.NoError(t, err)

		err = MaxLevelOption(0)(skiplist)
		testutil.ContainsError(t, err, fmt.Sprintf("attempted to set max level %d not in [1, 64]", 0))
	})

	t.Run("fail: max level > 64 is invalid", func(t *testing.T) {
		skiplist, err := New()
		require.NoError(t, err)

		err = MaxLevelOption(65)(skiplist)
		testutil.ContainsError(t, err, fmt.Sprintf("attempted to set max level %d not in [1, 64]", 65))
	})

	t.Run("pass: valid max level is set", func(t *testing.T) {
		skiplist, err := New()
		require.NoError(t, err)

		err = MaxLevelOption(32)(skiplist)
		require.NoError(t, err)

		assert.Equal(t, 32, skiplist.maxLevel)
	})
}

func TestProbabilityOption(t *testing.T) {
	t.Run("fail: non-skiplist is invalid", func(t *testing.T) {
		err := ProbabilityOption(0.25)(&rb.Tree{})
		testutil.ContainsError(t, err, "attempted to set probability on a non-skiplist")
	})

	t.Run("fail: probability <= 0 is invalid", func(t *testing.T) {
		skiplist, err := New()
		require.NoError(t, err)

		err = ProbabilityOption(0.)(skiplist)
		testutil.ContainsError(t, err, fmt.Sprintf("attempted to set probability %f not in (0, 1)", 0.))
	})

	t.Run("fail: probability >= 1 is invalid", func(t *testing.T) {
		skiplist, err := New()
		require.NoError(t, err)

		err = ProbabilityOption(1.)(skiplist)
		testutil.ContainsError(t, err, fmt.Sprintf("attempted to set probability %f not in (0, 1)", 1.))
	})

	t.Run("pass: valid probability is set", func(t *testing.T) {
		skiplist, err := New()
		require.NoError(t, err)

		err = ProbabilityOption(0.25)(skiplist)
		require.NoError(t, err)

		assert.Equal(t, 0.25, skiplist.p)
	})
}

func TestRandOption(t *testing.T) {
	t.Run("fail: non-skiplist is invalid", func(t *testing.T) {
		err := RandOption(rand.New(rand.NewSource(1)))(&rb.Tree{})
		testutil.ContainsError(t, err, "attempted to set rand source on a non-skiplist")
	})

	t.Run("pass: rand is set", func(t *testing.T) {
		skiplist, err := New()
		require.NoError(t, err)

		rand := rand.New(rand.NewSource(1))

		err = RandOption(rand)(skiplist)
		require.NoError(t, err)

		assert.Equal(t, rand, skiplist.rand)
	})
}
