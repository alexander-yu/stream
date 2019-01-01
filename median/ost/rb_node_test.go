package ost

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/alexander-yu/stream/util/testutil"
)

func TestColorString(t *testing.T) {
	c := Red
	assert.Equal(t, "Red", c.String())

	c = Black
	assert.Equal(t, "Black", c.String())
}

func TestRBNodeLeft(t *testing.T) {
	t.Run("pass: returns left child if node is not nil", func(t *testing.T) {
		node := NewRBNode(3)
		node.left = NewRBNode(4)
		left, err := node.Left()
		require.NoError(t, err)

		testutil.Approx(t, float64(4), left.Value())
	})

	t.Run("fail: return error if node is nil", func(t *testing.T) {
		var node *RBNode
		_, err := node.Left()
		assert.EqualError(t, err, "tried to retrieve child of nil node")
	})
}

func TestRBNodeRight(t *testing.T) {
	t.Run("pass: returns right child if node is not nil", func(t *testing.T) {
		node := NewRBNode(3)
		node.right = NewRBNode(4)
		right, err := node.Right()
		require.NoError(t, err)

		testutil.Approx(t, float64(4), right.Value())
	})

	t.Run("fail: return error if node is nil", func(t *testing.T) {
		var node *RBNode
		_, err := node.Right()
		assert.EqualError(t, err, "tried to retrieve child of nil node")
	})
}

func TestRBNodeColor(t *testing.T) {
	t.Run("pass: returns Black if node is nil", func(t *testing.T) {
		var node *RBNode
		assert.Equal(t, Black, node.Color())
	})

	t.Run("pass: returns color of node", func(t *testing.T) {
		node := NewRBNode(3)
		node = node.add(4)
		assert.Equal(t, Red, node.Color())
	})
}

func TestRBNodeSize(t *testing.T) {
	t.Run("pass: returns 0 if node is nil", func(t *testing.T) {
		var node *RBNode
		assert.Equal(t, 0, node.Size())
	})

	t.Run("pass: returns size of subtree", func(t *testing.T) {
		node := NewRBNode(3)
		node = node.add(4)
		assert.Equal(t, 2, node.Size())
	})
}

func TestRBNodeTreeString(t *testing.T) {
	t.Run("pass: returns empty string for empty tree", func(t *testing.T) {
		var node *RBNode
		assert.Equal(t, "", node.TreeString())
	})

	t.Run("pass: returns correct format for non-empty tree", func(t *testing.T) {
		var node *RBNode
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
				"        └── 1.000000",
				"            └── 1.000000",
				"",
			}, "\n"),
			node.TreeString(),
		)
	})
}
