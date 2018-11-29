package median

import (
	"fmt"
	"sync"

	"github.com/pkg/errors"
)

// Node represents a node in an order statistic tree.
type Node struct {
	left    *Node
	right   *Node
	val     float64
	_height int
	_size   int
	mux     sync.Mutex
}

func max(x int, y int) int {
	if x > y {
		return x
	}
	return y
}

// NewNode returns a new node containing a provided value.
func NewNode(val float64) *Node {
	return &Node{
		val:     val,
		_height: 0,
		_size:   1,
	}
}

// Left returns the left child of the node.
func (n *Node) Left() (*Node, error) {
	if n == nil {
		return nil, errors.New("tried to retrieve child of nil node")
	}
	n.mux.Lock()
	left := n.left
	n.mux.Unlock()
	return left, nil
}

// Right returns the right child of the node.
func (n *Node) Right() (*Node, error) {
	if n == nil {
		return nil, errors.New("tried to retrieve child of nil node")
	}
	n.mux.Lock()
	right := n.right
	n.mux.Unlock()
	return right, nil
}

// Height returns the height of the subtree rooted at the node.
func (n *Node) Height() int {
	// We need the extra nil check here because nil nodes do not have locks
	if n == nil {
		return -1
	}
	n.mux.Lock()
	height := n.height()
	n.mux.Unlock()
	return height
}

// Size returns the size of the subtree rooted at the node.
func (n *Node) Size() int {
	// We need the extra nil check here because nil nodes do not have locks
	if n == nil {
		return 0
	}
	n.mux.Lock()
	size := n.size()
	n.mux.Unlock()
	return size
}

// TreeString returns the string representation of the subtree rooted at the node.
func (n *Node) TreeString() string {
	if n == nil {
		return ""
	}
	n.mux.Lock()
	result := n.treeString("", "", true)
	n.mux.Unlock()
	return result
}

func (n *Node) height() int {
	if n == nil {
		return -1
	}
	return n._height
}

func (n *Node) size() int {
	if n == nil {
		return 0
	}
	return n._size
}

func (n *Node) add(val float64) *Node {
	if n == nil {
		return NewNode(val)
	} else if val < n.val {
		n.left = n.left.add(val)
	} else {
		n.right = n.right.add(val)
	}

	n._size = n.left.size() + n.right.size() + 1
	n._height = max(n.left.height(), n.right.height()) + 1
	return n.balance()
}

func (n *Node) remove(val float64) *Node {
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

	root._size = root.left.size() + root.right.size() + 1
	root._height = max(root.left.height(), root.right.height()) + 1
	return root.balance()
}

func (n *Node) min() *Node {
	if n.left == nil {
		return n
	}

	return n.left.min()
}

func (n *Node) removeMin() *Node {
	if n.left == nil {
		return n.right
	}

	n.left = n.left.removeMin()
	n._size = n.left.size() + n.right.size() + 1
	n._height = max(n.left.height(), n.right.height()) + 1
	return n.balance()
}

/*****************
 * Rotations
 *****************/

func (n *Node) balance() *Node {
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

func (n *Node) heightDiff() int {
	return n.left.height() - n.right.height()
}

func (n *Node) rotateLeft() *Node {
	m := n.right
	n.right = m.left
	m.left = n

	// No need to call size() here; we already know that n is not nil, since
	// rotations are only called for non-leaf nodes
	m._size = n._size
	n._size = n.left.size() + n.right.size() + 1

	n._height = max(n.left.height(), n.right.height()) + 1
	m._height = max(m.left.height(), m.right.height()) + 1

	return m
}

func (n *Node) rotateRight() *Node {
	m := n.left
	n.left = m.right
	m.right = n

	// No need to call size() here; we already know that n is not nil, since
	// rotations are only called for non-leaf nodes
	m._size = n._size
	n._size = n.left.size() + n.right.size() + 1

	n._height = max(n.left.height(), n.right.height()) + 1
	m._height = max(m.left.height(), m.right.height()) + 1

	return m
}

/*******************
 * Order Statistics
 *******************/

func (n *Node) get(i int) *Node {
	if n == nil {
		return nil
	}

	size := n.left.size()
	if i < size {
		return n.left.get(i)
	} else if i > size {
		return n.right.get(i - size - 1)
	}

	return n
}

func (n *Node) rank(val float64) int {
	if n == nil {
		return 0
	} else if val < n.val {
		return n.left.rank(val)
	} else if val > n.val {
		return 1 + n.left.size() + n.right.rank(val)
	}
	return n.left.size()
}

/*******************
 * Pretty-printing
 *******************/

func (n *Node) treeString(prefix string, result string, isTail bool) string {
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
