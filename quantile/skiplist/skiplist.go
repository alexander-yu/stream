package skiplist

import (
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/alexander-yu/stream/quantile/order"
	"github.com/pkg/errors"
)

const (
	// DefaultMaxLevel is the default max level for a skip list
	DefaultMaxLevel int = 12
	// DefaultProbability is the default probability for deciding levels in a skip list
	DefaultProbability float64 = 0.25
)

// Node represents a node in a skip list.
type Node struct {
	next  []*Node
	width []int
	val   float64
	head  bool
}

// Value returns the value stored at the node.
func (n *Node) Value() float64 {
	return n.val
}

func (n *Node) equals(m *Node) bool {
	return (n == nil && m == nil) || (n != nil && m != nil && n.val == m.val)
}

func (n *Node) string() string {
	if n == nil {
		return "nil"
	} else if n.head {
		return "head"
	}

	return strconv.FormatFloat(n.val, 'e', 9, 64)
}

// SkipList implements a skip list data structure,
// and also satisfies the order.Statistic interface.
type SkipList struct {
	head     *Node
	maxLevel int
	length   int
	rand     *rand.Rand
	p        float64
	probs    []float64
	prevs    []*Node
}

// New instantiates a SkipList struct.
func New(options ...order.Option) (*SkipList, error) {
	s := &SkipList{
		maxLevel: DefaultMaxLevel,
		p:        DefaultProbability,
		rand:     rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	for _, option := range options {
		err := option(s)
		if err != nil {
			return nil, errors.Wrap(err, "error setting option")
		}
	}

	// initialize head and tail nodes
	s.head = &Node{
		next:  make([]*Node, s.maxLevel),
		width: make([]int, s.maxLevel),
		head:  true,
	}
	for i := 0; i < s.maxLevel; i++ {
		s.head.width[i] = 1
	}
	s.prevs = make([]*Node, s.maxLevel)

	// generate table of probabilities of a new node having a given level
	for i := 1; i <= s.maxLevel; i++ {
		p := math.Pow(s.p, float64(i-1))
		s.probs = append(s.probs, p)
	}

	return s, nil
}

// Size returns the size of the skip list.
func (s *SkipList) Size() int {
	return s.length
}

// Clear resets the skip list.
func (s *SkipList) Clear() {
	s.head.next = make([]*Node, s.maxLevel)
	s.head.width = make([]int, s.maxLevel)
	for i := 0; i < s.maxLevel; i++ {
		s.head.width[i] = 1
	}
	s.prevs = make([]*Node, s.maxLevel)
	s.length = 0
}

// Add inserts a value into the skip list.
func (s *SkipList) Add(val float64) {
	prevs := s.getPrevs(val)
	level := s.randLevel()
	node := &Node{
		next:  make([]*Node, level),
		width: make([]int, level),
		val:   val,
	}

	for i := 0; i < level; i++ {
		node.next[i] = prevs[i].next[i]
		prevs[i].next[i] = node
	}

	// update widths
	node.width[0] = 1
	for i := 1; i < s.maxLevel; i++ {
		prevs[i].width[i]++
	}
	for i := 1; i < level; i++ {
		for curr := node; !curr.equals(node.next[i]); curr = curr.next[i-1] {
			node.width[i] += curr.width[i-1]
		}
		prevs[i].width[i] -= node.width[i]
	}

	s.length++
}

// Remove deletes a value from the skip list.
func (s *SkipList) Remove(val float64) {
	prevs := s.getPrevs(val)
	// if node with value is found, then set all predecessors to point
	// to the nodes in node.next, and update widths
	if node := prevs[0].next[0]; node != nil && node.val == val {
		for i := range prevs {
			if i < len(node.next) {
				prevs[i].next[i] = node.next[i]
				prevs[i].width[i] += node.width[i] - 1
			} else {
				prevs[i].width[i]--
			}
		}
		s.length--
	} else {

	}
}

// Select returns the node with the kth smallest value in the skip list.
func (s *SkipList) Select(k int) order.Node {
	if k < 0 || k >= s.length {
		return nil
	}

	node := s.head
	pos := 0
	k++ // increment k because the head node is counted
	for i := s.maxLevel - 1; i >= 0; i-- {
		for pos+node.width[i] <= k {
			pos += node.width[i]
			node = node.next[i]
		}
	}

	return node
}

// Rank returns the number of nodes strictly less than the given value.
func (s *SkipList) Rank(val float64) int {
	rank := 0
	node := s.head
	for i := s.maxLevel - 1; i >= 0; i-- {
		for node.next[i] != nil && node.next[i].val < val {
			rank += node.width[i]
			node = node.next[i]
		}
	}

	return rank
}

// String returns the string representation of the skip list.
func (s *SkipList) String() string {
	result := ""
	for i := s.maxLevel - 1; i >= 0; i-- {
		for prev := s.head; prev != nil; prev = prev.next[i] {
			result += prev.string() + strings.Repeat("-", prev.width[i])
		}
		result += "tail\n"
	}
	return result
}

// getPrevs retrieves the list of nodes at each level of the skip list
// that would be the immediate predecessors of the given value.
func (s *SkipList) getPrevs(val float64) []*Node {
	prev := s.head

	// walk along the list starting from the highest level and
	// descend down to more granular levels once you overshoot on the
	// the current level
	for i := s.maxLevel - 1; i >= 0; i-- {
		next := prev.next[i]
		for next != nil && next.val < val {
			prev = next
			next = next.next[i]
		}
		s.prevs[i] = prev
	}

	return s.prevs
}

// randLevel generates at random the max level that a node will be in
func (s *SkipList) randLevel() int {
	r := s.rand.Float64()
	level := 1
	for level < s.maxLevel && r < s.probs[level] {
		level++
	}
	return level
}
