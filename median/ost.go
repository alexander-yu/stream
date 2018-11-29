package median

// OrderStatisticTree is a variant of the binary tree that provides two additional
// methods: Select(i), which finds the ith smallest element in the tree, and Rank(x),
// which finds the rank of element x in the tree. This is implemented as an AVL tree
// whose nodes keep track of the sizes of their subtrees.
type OrderStatisticTree struct {
	root *Node
}

// Size returns the size of the tree.
func (t *OrderStatisticTree) Size() int {
	return t.root.Size()
}

// Height returns the height of the tree.
func (t *OrderStatisticTree) Height() int {
	return t.root.Height()
}

// Add inserts a value into the tree.
func (t *OrderStatisticTree) Add(val float64) {
	t.root = t.root.add(val)
}

// Remove deletes a value from the tree.
func (t *OrderStatisticTree) Remove(val float64) {
	t.root = t.root.remove(val)
}

// Select returns the node with the ith smallest value in the tree.
func (t *OrderStatisticTree) Select(i int) *Node {
	return t.root.get(i)
}

// Rank returns the number of nodes strictly less than the value.
func (t *OrderStatisticTree) Rank(val float64) int {
	return t.root.rank(val)
}

// String returns the string representation of the tree.
func (t *OrderStatisticTree) String() string {
	return t.root.treeString("", "", true)
}
