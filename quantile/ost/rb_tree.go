package ost

import "github.com/alexander-yu/stream/quantile/order"

// RBTree implements a red-black tree data structure,
// and also satisfies the Tree interface.
type RBTree struct {
	root *RBNode
}

// Size returns the size of the tree.
func (t *RBTree) Size() int {
	return t.root.Size()
}

// Add inserts a value into the tree.
func (t *RBTree) Add(val float64) {
	t.root = t.root.add(val)
}

// Remove deletes a value from the tree.
func (t *RBTree) Remove(val float64) {
	t.root = t.root.remove(val)
}

// Select returns the node with the ith smallest value in the tree.
func (t *RBTree) Select(i int) order.Node {
	return t.root.Select(i)
}

// Rank returns the number of nodes strictly less than the value.
func (t *RBTree) Rank(val float64) int {
	return t.root.Rank(val)
}

// String returns the string representation of the tree.
func (t *RBTree) String() string {
	return t.root.TreeString()
}

// Clear resets the tree.
func (t *RBTree) Clear() {
	*t = RBTree{}
}
