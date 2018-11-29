package median

import (
	"sync"
)

// OrderStatisticTree is a variant of the binary tree that provides two additional
// methods: Select(i), which finds the ith smallest element in the tree, and Rank(x),
// which finds the rank of element x in the tree. This is implemented as an AVL tree
// whose nodes keep track of the sizes of their subtrees.
type OrderStatisticTree struct {
	root *Node
	mux  sync.Mutex
}

// Size returns the size of the tree.
func (t *OrderStatisticTree) Size() int {
	t.mux.Lock()
	size := t.root.size()
	t.mux.Unlock()
	return size
}

// Height returns the height of the tree.
func (t *OrderStatisticTree) Height() int {
	t.mux.Lock()
	height := t.root.height()
	t.mux.Unlock()
	return height
}

// Add inserts a value into the tree.
func (t *OrderStatisticTree) Add(val float64) {
	t.mux.Lock()
	t.root = t.root.add(val)
	t.mux.Unlock()
}

// Remove deletes a value from the tree.
func (t *OrderStatisticTree) Remove(val float64) {
	t.mux.Lock()
	t.root = t.root.remove(val)
	t.mux.Unlock()
}

// Select returns the node with the ith smallest value in the tree.
func (t *OrderStatisticTree) Select(i int) *Node {
	t.mux.Lock()
	node := t.root.get(i)
	t.mux.Unlock()
	return node
}

// Rank returns the number of nodes strictly less than the value.
func (t *OrderStatisticTree) Rank(val float64) int {
	t.mux.Lock()
	rank := t.root.rank(val)
	t.mux.Unlock()
	return rank
}

// String returns the string representation of the tree.
func (t *OrderStatisticTree) String() string {
	t.mux.Lock()
	result := t.root.treeString("", "", true)
	t.mux.Unlock()
	return result
}
