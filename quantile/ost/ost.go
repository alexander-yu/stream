package ost

import (
	"github.com/alexander-yu/stream/quantile/order"
)

// Tree is the interface for order statistic trees, which are variants of binary trees
// that provide two additional methods: Select(i), which finds the ith smallest element
// in the tree, and Rank(x), which finds the rank of element x in the tree.
type Tree interface {
	Size() int
	Add(float64)
	Remove(float64)
	Select(int) order.Node
	Rank(float64) int
	String() string
	Clear()
}

// Node is the interface for any node struct within an Tree.
type Node interface {
	Left() (Node, error)
	Right() (Node, error)
	Size() int
	Value() float64
	Select(int) order.Node
	Rank(float64) int
	TreeString() string
}
