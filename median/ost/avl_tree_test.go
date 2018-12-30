package ost

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testData() *AVLTree {
	tree := &AVLTree{}
	tree.Add(5)
	tree.Add(6)
	tree.Add(7)
	tree.Add(3)
	tree.Add(4)
	tree.Add(1)
	tree.Add(2)
	tree.Add(1)
	return tree
}

func TestAVLTreeAdd(t *testing.T) {
	tree := testData()

	assert.Equal(t, 8, tree.Size())
	assert.Equal(t, 3, tree.Height())
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
		tree.String(),
	)

	tree.Add(6.5)
	tree.Add(6.75)
	tree.Add(6.25)
	assert.Equal(t, 11, tree.Size())
	assert.Equal(t, 3, tree.Height())
	assert.Equal(
		t,
		strings.Join([]string{
			"│           ┌── 7.000000",
			"│       ┌── 6.750000",
			"│   ┌── 6.500000",
			"│   │   │   ┌── 6.250000",
			"│   │   └── 6.000000",
			"│   │       └── 5.000000",
			"└── 4.000000",
			"    │   ┌── 3.000000",
			"    └── 2.000000",
			"        │   ┌── 1.000000",
			"        └── 1.000000",
			"",
		}, "\n"),
		tree.String(),
	)
}

func TestAVLTreeRemove(t *testing.T) {
	tree := testData()
	tree.Remove(5)
	tree.Remove(7)

	assert.Equal(t, 6, tree.Size())
	assert.Equal(t, 2, tree.Height())
	assert.Equal(
		t,
		strings.Join([]string{
			"│       ┌── 6.000000",
			"│   ┌── 4.000000",
			"│   │   └── 3.000000",
			"└── 2.000000",
			"    │   ┌── 1.000000",
			"    └── 1.000000",
			"",
		}, "\n"),
		tree.String(),
	)

	tree.Remove(2)
	assert.Equal(t, 5, tree.Size())
	assert.Equal(t, 2, tree.Height())
	assert.Equal(
		t,
		strings.Join([]string{
			"│       ┌── 6.000000",
			"│   ┌── 4.000000",
			"└── 3.000000",
			"    │   ┌── 1.000000",
			"    └── 1.000000",
			"",
		}, "\n"),
		tree.String(),
	)

	tree.Remove(1)
	tree.Remove(6)
	tree.Remove(4)
	tree.Remove(3)
	assert.Equal(t, 1, tree.Size())
	assert.Equal(t, 0, tree.Height())
	assert.Equal(
		t,
		strings.Join([]string{
			"└── 1.000000",
			"",
		}, "\n"),
		tree.String(),
	)
}

func TestAVLTreeRank(t *testing.T) {
	tree := testData()
	rank := tree.Rank(3)
	assert.Equal(t, 3, rank)

	rank = tree.Rank(5.5)
	assert.Equal(t, 6, rank)

	rank = tree.Rank(-1)
	assert.Equal(t, 0, rank)
}

func TestAVLTreeSelect(t *testing.T) {
	tree := testData()
	node := tree.Select(5)
	assert.Equal(t, float64(5), node.Value())

	node = tree.Select(-1)
	assert.Nil(t, node)

	node = tree.Select(9)
	assert.Nil(t, node)
}
