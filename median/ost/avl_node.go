package ost

import (
	"fmt"

	"github.com/pkg/errors"
)

// AVLNode represents a node in an AVL tree.
type AVLNode struct {
	left   *AVLNode
	right  *AVLNode
	val    float64
	height int
	size   int
}

func max(x int, y int) int {
	if x > y {
		return x
	}
	return y
}

// NewAVLNode instantiates a AVLNode struct with a a provided value.
func NewAVLNode(val float64) *AVLNode {
	return &AVLNode{
		val:    val,
		height: 0,
		size:   1,
	}
}

// Left returns the left child of the node.
func (n *AVLNode) Left() (Node, error) {
	if n == nil {
		return nil, errors.New("tried to retrieve child of nil node")
	}
	return n.left, nil
}

// Right returns the right child of the node.
func (n *AVLNode) Right() (Node, error) {
	if n == nil {
		return nil, errors.New("tried to retrieve child of nil node")
	}
	return n.right, nil
}

// Height returns the height of the subtree rooted at the node.
func (n *AVLNode) Height() int {
	if n == nil {
		return -1
	}
	return n.height
}

// Size returns the size of the subtree rooted at the node.
func (n *AVLNode) Size() int {
	if n == nil {
		return 0
	}
	return n.size
}

// Value returns the value stored at the node.
func (n *AVLNode) Value() float64 {
	return n.val
}

// TreeString returns the string representation of the subtree rooted at the node.
func (n *AVLNode) TreeString() string {
	if n == nil {
		return ""
	}
	return n.treeString("", "", true)
}

func (n *AVLNode) add(val float64) *AVLNode {
	if n == nil {
		return NewAVLNode(val)
	} else if val <= n.val {
		n.left = n.left.add(val)
	} else {
		n.right = n.right.add(val)
	}

	n.size = n.left.Size() + n.right.Size() + 1
	n.height = max(n.left.Height(), n.right.Height()) + 1
	return n.balance()
}

func (n *AVLNode) remove(val float64) *AVLNode {
	root := n
	if val < root.val {
		root.left = root.left.remove(val)
	} else if val > root.val {
		root.right = root.right.remove(val)
	} else {
		if root.left == nil {
			return root.right
		} else if root.right == nil {
			return root.left
		}
		root = n.right.min()
		root.right = n.right.removeMin()
		root.left = n.left
	}

	root.size = root.left.Size() + root.right.Size() + 1
	root.height = max(root.left.Height(), root.right.Height()) + 1
	return root.balance()
}

func (n *AVLNode) min() *AVLNode {
	if n.left == nil {
		return n
	}

	return n.left.min()
}

func (n *AVLNode) removeMin() *AVLNode {
	if n.left == nil {
		return n.right
	}

	n.left = n.left.removeMin()
	n.size = n.left.Size() + n.right.Size() + 1
	n.height = max(n.left.Height(), n.right.Height()) + 1
	return n.balance()
}

/*****************
 * Rotations
 *****************/

func (n *AVLNode) balance() *AVLNode {
	if n.heightDiff() < -1 {
		// Since we've entered this block, we already
		// know that the right child is not nil
		if n.right.heightDiff() > 0 {
			n.right = n.right.rotateRight()
		}
		return n.rotateLeft()
	} else if n.heightDiff() > 1 {
		// Since we've entered this block, we already
		// know that the left child is not nil
		if n.left.heightDiff() < 0 {
			n.left = n.left.rotateLeft()
		}
		return n.rotateRight()
	}

	return n
}

func (n *AVLNode) heightDiff() int {
	return n.left.Height() - n.right.Height()
}

func (n *AVLNode) rotateLeft() *AVLNode {
	m := n.right
	n.right = m.left
	m.left = n

	// No need to call size() here; we already know that n is not nil, since
	// rotations are only called for non-leaf nodes
	m.size = n.size
	n.size = n.left.Size() + n.right.Size() + 1

	n.height = max(n.left.Height(), n.right.Height()) + 1
	m.height = max(m.left.Height(), m.right.Height()) + 1

	return m
}

func (n *AVLNode) rotateRight() *AVLNode {
	m := n.left
	n.left = m.right
	m.right = n

	// No need to call size() here; we already know that n is not nil, since
	// rotations are only called for non-leaf nodes
	m.size = n.size
	n.size = n.left.Size() + n.right.Size() + 1

	n.height = max(n.left.Height(), n.right.Height()) + 1
	m.height = max(m.left.Height(), m.right.Height()) + 1

	return m
}

/*******************
 * Order Statistics
 *******************/

// Select returns the node with the ith smallest value in the
// subtree rooted at the node..
func (n *AVLNode) Select(i int) Node {
	if n == nil {
		return nil
	}

	size := n.left.Size()
	if i < size {
		return n.left.Select(i)
	} else if i > size {
		return n.right.Select(i - size - 1)
	}

	return n
}

// Rank returns the number of nodes strictly less than the value that
// are contained in the subtree rooted at the node.
func (n *AVLNode) Rank(val float64) int {
	if n == nil {
		return 0
	} else if val < n.val {
		return n.left.Rank(val)
	} else if val > n.val {
		return 1 + n.left.Size() + n.right.Rank(val)
	}
	return n.left.Size()
}

/*******************
 * Pretty-printing
 *******************/

func (n *AVLNode) treeString(prefix string, result string, isTail bool) string {
	if n.right != nil {
		if isTail {
			result = n.right.treeString(fmt.Sprintf("%s│   ", prefix), result, false)
		} else {
			result = n.right.treeString(fmt.Sprintf("%s    ", prefix), result, false)
		}
	}

	if isTail {
		result = fmt.Sprintf("%s%s└── %f\n", result, prefix, n.val)
	} else {
		result = fmt.Sprintf("%s%s┌── %f\n", result, prefix, n.val)
	}

	if n.left != nil {
		if isTail {
			result = n.left.treeString(fmt.Sprintf("%s    ", prefix), result, true)
		} else {
			result = n.left.treeString(fmt.Sprintf("%s│   ", prefix), result, true)
		}
	}

	return result
}
