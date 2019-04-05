package quantile

import (
	"fmt"
	"math"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestImplInit(t *testing.T) {
	t.Run("pass: AVL implementation is supported", func(t *testing.T) {
		i := AVL
		_, err := i.init()
		assert.NoError(t, err)
	})

	t.Run("pass: red black implementation is supported", func(t *testing.T) {
		i := RedBlack
		_, err := i.init()
		assert.NoError(t, err)
	})

	t.Run("pass: skip list implementation is supported", func(t *testing.T) {
		i := SkipList
		_, err := i.init()
		assert.NoError(t, err)
	})

	t.Run("fail: unsupported implementations return an error", func(t *testing.T) {
		i := Impl(-1)
		_, err := i.init()
		assert.EqualError(t, err, fmt.Sprintf("%v is not a supported Impl value", i))
	})
}

func BenchmarkImplAdd(b *testing.B) {
	for k := 3.; k < 20; k++ {
		n := int(math.Pow(2, k))
		xs := make([]float64, n)
		for i := 0; i < n; i++ {
			xs[i] = 200*rand.Float64() - 100
		}
		y := 200*rand.Float64() - 100

		avl, err := AVL.init()
		require.NoError(b, err)

		rb, err := RedBlack.init()
		require.NoError(b, err)

		skiplist, err := SkipList.init()
		require.NoError(b, err)

		for _, x := range xs {
			avl.Add(x)
			rb.Add(x)
			skiplist.Add(x)
		}

		b.Run(fmt.Sprintf("AVL tree [%d]", n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StartTimer()
				avl.Add(y)
				b.StopTimer()
				avl.Remove(y)
			}
		})

		b.Run(fmt.Sprintf("red-black tree [%d]", n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StartTimer()
				rb.Add(y)
				b.StopTimer()
				rb.Remove(y)
			}
		})

		b.Run(fmt.Sprintf("skip list [%d]", n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StartTimer()
				skiplist.Add(y)
				b.StopTimer()
				skiplist.Remove(y)
			}
		})
	}
}

func BenchmarkImplRemove(b *testing.B) {
	for k := 3.; k < 20; k++ {
		n := int(math.Pow(2, k))
		xs := make([]float64, n)
		for i := 0; i < n; i++ {
			xs[i] = 200*rand.Float64() - 100
		}
		y := 200*rand.Float64() - 100

		avl, err := AVL.init()
		require.NoError(b, err)

		rb, err := RedBlack.init()
		require.NoError(b, err)

		skiplist, err := SkipList.init()
		require.NoError(b, err)

		for _, x := range xs {
			avl.Add(x)
			rb.Add(x)
			skiplist.Add(x)
		}

		b.Run(fmt.Sprintf("AVL tree [%d]", n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StartTimer()
				avl.Remove(y)
				b.StopTimer()
				avl.Add(y)
			}
		})

		b.Run(fmt.Sprintf("red-black tree [%d]", n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StartTimer()
				rb.Remove(y)
				b.StopTimer()
				rb.Add(y)
			}
		})

		b.Run(fmt.Sprintf("skip list [%d]", n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StartTimer()
				skiplist.Remove(y)
				b.StopTimer()
				skiplist.Add(y)
			}
		})
	}
}

func BenchmarkImplSelect(b *testing.B) {
	for k := 3.; k < 20; k++ {
		n := int(math.Pow(2, k))
		xs := make([]float64, n)
		for i := 0; i < n; i++ {
			xs[i] = 200*rand.Float64() - 100
		}
		idx := rand.Intn(n)

		avl, err := AVL.init()
		require.NoError(b, err)

		rb, err := RedBlack.init()
		require.NoError(b, err)

		skiplist, err := SkipList.init()
		require.NoError(b, err)

		for _, x := range xs {
			avl.Add(x)
			rb.Add(x)
			skiplist.Add(x)
		}

		b.Run(fmt.Sprintf("AVL tree [%d-%d]", n, idx), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				avl.Select(idx)
			}
		})

		b.Run(fmt.Sprintf("red-black tree [%d-%d]", n, idx), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				rb.Select(idx)
			}
		})

		b.Run(fmt.Sprintf("skip list [%d-%d]", n, idx), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				skiplist.Select(idx)
			}
		})
	}
}

func BenchmarkImplRank(b *testing.B) {
	for k := 3.; k < 20; k++ {
		n := int(math.Pow(2, k))
		xs := make([]float64, n)
		for i := 0; i < n; i++ {
			xs[i] = 200*rand.Float64() - 100
		}
		y := 200*rand.Float64() - 100

		avl, err := AVL.init()
		require.NoError(b, err)

		rb, err := RedBlack.init()
		require.NoError(b, err)

		skiplist, err := SkipList.init()
		require.NoError(b, err)

		for _, x := range xs {
			avl.Add(x)
			rb.Add(x)
			skiplist.Add(x)
		}

		b.Run(fmt.Sprintf("AVL tree [%d]", n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				avl.Rank(y)
			}
		})

		b.Run(fmt.Sprintf("red-black tree [%d]", n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				rb.Rank(y)
			}
		})

		b.Run(fmt.Sprintf("skip list [%d]", n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				skiplist.Rank(y)
			}
		})
	}
}
