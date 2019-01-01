package ost

import "github.com/pkg/errors"

// Tree is the interface for order statistic trees, which are variants of binary trees
// that provide two additional methods: Select(i), which finds the ith smallest element
// in the tree, and Rank(x), which finds the rank of element x in the tree.
type Tree interface {
	Size() int
	Add(val float64)
	Remove(val float64)
	Select(int) Node
	Rank(float64) int
	String() string
}

// Node is the interface for any node struct within an Tree.
type Node interface {
	Left() (Node, error)
	Right() (Node, error)
	Size() int
	Value() float64
	Select(int) Node
	Rank(float64) int
	TreeString() string
}

// Impl represents an enum that enumerates the currently supported implementations
// for the Tree interface.
type Impl int

const (
	// AVL represents the AVL tree implementation for the Tree interface
	AVL Impl = 1
	// RB represents the red black tree implementation for the Tree interface
	RB Impl = 2
)

// EmptyTree returns an empty Tree struct, depending on which implementation
// is being called.
func (i Impl) EmptyTree() (Tree, error) {
	switch i {
	case AVL:
		return &AVLTree{}, nil
	case RB:
		return &RBTree{}, nil
	default:
		return nil, errors.Errorf("%v is not a supported OST implementation", i)
	}
}
