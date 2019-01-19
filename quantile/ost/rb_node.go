package ost

import (
	"fmt"

	"github.com/pkg/errors"
)

// Color represents the color of the node.
type Color bool

// The only allowed colors are red and black.
const (
	Red   Color = true
	Black Color = false
)

func (c Color) String() string {
	switch c {
	case Black:
		return "Black"
	default:
		return "Red"
	}
}

// RBNode represents a node in a red black tree.
type RBNode struct {
	left  *RBNode
	right *RBNode
	val   float64
	color Color
	size  int
}

// NewRBNode instantiates a RBNode struct with a a provided value.
func NewRBNode(val float64) *RBNode {
	return &RBNode{
		val:   val,
		color: Red,
		size:  1,
	}
}

// Left returns the left child of the node.
func (n *RBNode) Left() (Node, error) {
	if n == nil {
		return nil, errors.New("tried to retrieve child of nil node")
	}
	return n.left, nil
}

// Right returns the right child of the node.
func (n *RBNode) Right() (Node, error) {
	if n == nil {
		return nil, errors.New("tried to retrieve child of nil node")
	}
	return n.right, nil
}

// Size returns the size of the subtree rooted at the node.
func (n *RBNode) Size() int {
	if n == nil {
		return 0
	}
	return n.size
}

// Value returns the value stored at the node.
func (n *RBNode) Value() float64 {
	return n.val
}

// Color returns the color of the node.
// By default, nil nodes are black.
func (n *RBNode) Color() Color {
	if n == nil {
		return Black
	}
	return n.color
}

// TreeString returns the string representation of the subtree rooted at the node.
func (n *RBNode) TreeString() string {
	if n == nil {
		return ""
	}
	return n.treeString("", "", true)
}

func (n *RBNode) add(val float64) *RBNode {
	if n == nil {
		return NewRBNode(val)
	} else if val <= n.val {
		n.left = n.left.add(val)
	} else {
		n.right = n.right.add(val)
	}
	return n.addBalance()
}

func (n *RBNode) remove(val float64) *RBNode {
	if !n.contains(val) {
		return n
	}

	if val < n.val {
		if n.left.Color() == Black && n.left.left.Color() == Black {
			n = n.moveRedLeft()
		}
		n.left = n.left.remove(val)
	} else {
		if n.left.Color() == Red {
			n = n.rotateRight()
		}
		if val == n.val && n.right == nil {
			return nil
		}
		if n.right.Color() == Black && n.right.left.Color() == Black {
			n = n.moveRedRight()
		}
		if val == n.val {
			x := n.right.min()
			n.val = x.val
			n.right = n.right.removeMin()
		} else {
			n.right = n.right.remove(val)
		}
	}

	return n.removeBalance()
}

func (n *RBNode) removeMin() *RBNode {
	if n.left == nil {
		return nil
	}
	if n.left.Color() == Black && n.left.left.Color() == Black {
		n = n.moveRedLeft()
	}

	n.left = n.left.removeMin()
	return n.removeBalance()
}

func (n *RBNode) min() *RBNode {
	if n.left == nil {
		return n
	}
	return n.left.min()
}

func (n *RBNode) contains(val float64) bool {
	for n != nil {
		if val == n.val {
			return true
		} else if val < n.val {
			n = n.left
		} else {
			n = n.right
		}
	}
	return false
}

/*****************
 * Rotations
 *****************/

func (n *RBNode) addBalance() *RBNode {
	if n.left.Color() == Black && n.right.Color() == Red {
		n = n.rotateLeft()
	}
	if n.left.Color() == Red && n.left.left.Color() == Red {
		n = n.rotateRight()
	}
	if n.left.Color() == Red && n.right.Color() == Red {
		n.flipColors()
	}

	n.size = n.left.Size() + n.right.Size() + 1
	return n
}

func (n *RBNode) removeBalance() *RBNode {
	if n.right.Color() == Red {
		n = n.rotateLeft()
	}
	if n.left.Color() == Red && n.left.left.Color() == Red {
		n = n.rotateRight()
	}
	if n.left.Color() == Red && n.right.Color() == Red {
		n.flipColors()
	}

	n.size = n.left.Size() + n.right.Size() + 1
	return n
}

func (n *RBNode) rotateLeft() *RBNode {
	x := n.right
	n.right = x.left
	x.left = n
	x.color = x.left.color
	x.left.color = Red
	x.size = n.size
	n.size = n.left.Size() + n.right.Size() + 1
	return x
}

func (n *RBNode) rotateRight() *RBNode {
	x := n.left
	n.left = x.right
	x.right = n
	x.color = x.right.color
	x.right.color = Red
	x.size = n.size
	n.size = n.left.Size() + n.right.Size() + 1
	return x
}

func (n *RBNode) flipColors() {
	n.color = !n.color
	n.left.color = !n.left.color
	n.right.color = !n.right.color
}

func (n *RBNode) moveRedLeft() *RBNode {
	n.flipColors()
	if n.right.left.Color() == Red {
		n.right = n.right.rotateRight()
		n = n.rotateLeft()
		n.flipColors()
	}
	return n
}

func (n *RBNode) moveRedRight() *RBNode {
	n.flipColors()
	if n.left.left.Color() == Red {
		n = n.rotateRight()
		n.flipColors()
	}
	return n
}

/*******************
 * Order Statistics
 *******************/

// Select returns the node with the ith smallest value in the
// subtree rooted at the node..
func (n *RBNode) Select(i int) Node {
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
func (n *RBNode) Rank(val float64) int {
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

// treeString recursively prints out a subtree rooted at the node in a sideways format, as below:
// │       ┌── 7.000000
// │   ┌── 6.000000
// │   │   └── 5.000000
// └── 4.000000
//     │   ┌── 3.000000
//     └── 2.000000
//         └── 1.000000
//             └── 1.000000
func (n *RBNode) treeString(prefix string, result string, isTail bool) string {
	// isTail indicates whether or not the current node's parent branch needs to be represented
	// as a "tail", i.e. its branch needs to hang in the string representation, rather than branch upwards.
	if isTail {
		// If true, then we need to print the subtree like this:
		// │   ┌── [n.right.treeString()]
		// └── [n.val]
		//     └── [n.left.treeString()]
		if n.right != nil {
			result = n.right.treeString(fmt.Sprintf("%s│   ", prefix), result, false)
		}
		result = fmt.Sprintf("%s%s└── %f\n", result, prefix, n.val)
		if n.left != nil {
			result = n.left.treeString(fmt.Sprintf("%s    ", prefix), result, true)
		}
	} else {
		// If false, then we need to print the subtree like this:
		//     ┌── [n.right.treeString()]
		// ┌── [n.val]
		// │   └── [n.left.treeString()]
		if n.right != nil {
			result = n.right.treeString(fmt.Sprintf("%s    ", prefix), result, false)
		}
		result = fmt.Sprintf("%s%s┌── %f\n", result, prefix, n.val)
		if n.left != nil {
			result = n.left.treeString(fmt.Sprintf("%s│   ", prefix), result, true)
		}
	}

	return result
}
