package ost

// AVLTree implements an AVL tree data structure,
// and also satisfies the Tree interface.
type AVLTree struct {
	root *AVLNode
}

// Size returns the size of the tree.
func (t *AVLTree) Size() int {
	return t.root.Size()
}

// Height returns the height of the tree.
func (t *AVLTree) Height() int {
	return t.root.Height()
}

// Add inserts a value into the tree.
func (t *AVLTree) Add(val float64) {
	t.root = t.root.add(val)
}

// Remove deletes a value from the tree.
func (t *AVLTree) Remove(val float64) {
	t.root = t.root.remove(val)
}

// Select returns the node with the ith smallest value in the tree.
func (t *AVLTree) Select(i int) Node {
	return t.root.Select(i)
}

// Rank returns the number of nodes strictly less than the value.
func (t *AVLTree) Rank(val float64) int {
	return t.root.Rank(val)
}

// String returns the string representation of the tree.
func (t *AVLTree) String() string {
	return t.root.TreeString()
}
