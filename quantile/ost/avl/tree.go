package avl

import "github.com/alexander-yu/stream/quantile/order"

// Tree implements an AVL tree data structure,
// and also satisfies the ost.Tree interface,
// as well as the order.Statistic interface.
type Tree struct {
	root *Node
}

// Size returns the size of the tree.
func (t *Tree) Size() int {
	return t.root.Size()
}

// Height returns the height of the tree.
func (t *Tree) Height() int {
	return t.root.Height()
}

// Add inserts a value into the tree.
func (t *Tree) Add(val float64) {
	t.root = t.root.add(val)
}

// Remove deletes a value from the tree.
func (t *Tree) Remove(val float64) {
	t.root = t.root.remove(val)
}

// Select returns the node with the ith smallest value in the tree.
func (t *Tree) Select(i int) order.Node {
	return t.root.Select(i)
}

// Rank returns the number of nodes strictly less than the value.
func (t *Tree) Rank(val float64) int {
	return t.root.Rank(val)
}

// String returns the string representation of the tree.
func (t *Tree) String() string {
	return t.root.TreeString()
}

// Clear resets the tree.
func (t *Tree) Clear() {
	*t = Tree{}
}
