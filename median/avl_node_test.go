package median

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/alexander-yu/stream/util/testutil"
)

func TestAVLNodeLeft(t *testing.T) {
	t.Run("pass: returns left child if node is not nil", func(t *testing.T) {
		node := NewAVLNode(3)
		node.left = NewAVLNode(4)
		left, err := node.Left()
		require.NoError(t, err)

		testutil.Approx(t, float64(4), left.val)
	})

	t.Run("fail: return error if node is nil", func(t *testing.T) {
		var node *AVLNode
		_, err := node.Left()
		assert.EqualError(t, err, "tried to retrieve child of nil node")
	})
}

func TestAVLNodeRight(t *testing.T) {
	t.Run("pass: returns right child if node is not nil", func(t *testing.T) {
		node := NewAVLNode(3)
		node.right = NewAVLNode(4)
		right, err := node.Right()
		require.NoError(t, err)

		testutil.Approx(t, float64(4), right.val)
	})

	t.Run("fail: return error if node is nil", func(t *testing.T) {
		var node *AVLNode
		_, err := node.Right()
		assert.EqualError(t, err, "tried to retrieve child of nil node")
	})
}

func TestAVLNodeHeight(t *testing.T) {
	t.Run("pass: returns -1 if node is nil", func(t *testing.T) {
		var node *AVLNode
		assert.Equal(t, -1, node.Height())
	})

	t.Run("pass: returns height of subtree", func(t *testing.T) {
		node := NewAVLNode(3)
		node = node.add(4)
		assert.Equal(t, 1, node.Height())
	})
}

func TestAVLNodeSize(t *testing.T) {
	t.Run("pass: returns 0 if node is nil", func(t *testing.T) {
		var node *AVLNode
		assert.Equal(t, 0, node.Size())
	})

	t.Run("pass: returns size of subtree", func(t *testing.T) {
		node := NewAVLNode(3)
		node = node.add(4)
		assert.Equal(t, 2, node.Size())
	})
}

func TestAVLNodeTreeString(t *testing.T) {
	t.Run("pass: returns empty string for empty tree", func(t *testing.T) {
		var node *AVLNode
		assert.Equal(t, "", node.TreeString())
	})

	t.Run("pass: returns correct format for non-empty tree", func(t *testing.T) {
		var node *AVLNode
		node = node.add(5)
		node = node.add(6)
		node = node.add(7)
		node = node.add(3)
		node = node.add(4)
		node = node.add(1)
		node = node.add(2)
		node = node.add(1)
		assert.Equal(
			t,
			strings.Join([]string{
				"│       ┌── 7.000000",
				"│   ┌── 6.000000",
				"│   │   └── 5.000000",
				"└── 4.000000",
				"    │   ┌── 3.000000",
				"    └── 2.000000",
				"        │   ┌── 1.000000",
				"        └── 1.000000",
				"",
			}, "\n"),
			node.TreeString(),
		)
	})
}
