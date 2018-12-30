package median

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
}

// Node is the interface for any node struct within an OrderStatisticTree.
type Node interface {
	Height() int
	Size() int
	Value() float64
	get(int) Node
	rank(float64) int
}
