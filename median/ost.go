package median

import "github.com/pkg/errors"

// OrderStatisticTree is a variant of the binary tree that provides two additional
// methods: Select(i), which finds the ith smallest element in the tree, and Rank(x),
// which finds the rank of element x in the tree.
type OrderStatisticTree interface {
	Size() int
	Height() int
	Add(val float64)
	Remove(val float64)
	Select(int) Node
	Rank(float64) int
	String() string
}

// Node is the interface for any node struct within an OrderStatisticTree.
type Node interface {
	Height() int
	Size() int
	Value() float64
	get(int) Node
	rank(float64) int
}

// OSTImpl represents an enum that enumerates the currently supported implementations
// for the OrderStatisticTree interface.
type OSTImpl int

const (
	// AVL represents the AVL tree implementation for the OrderStatisticTree interface
	AVL OSTImpl = 1
)

// EmptyTree returns an empty OrderStatisticTree struct, depending on which implementation
// is being called.
func (i OSTImpl) EmptyTree() (OrderStatisticTree, error) {
	switch i {
	case AVL:
		return &AVLTree{}, nil
	default:
		return nil, errors.Errorf("%v is not a supported OST implementation", i)
	}
}
