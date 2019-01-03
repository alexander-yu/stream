package heap

import heapops "container/heap"

// Heap implements a heap data structure.
type Heap struct {
	items []*Item
	cmp   func(float64, float64) bool
}

// Item represents an item in the heap; in particular,
// it contains not only the value, but also the index
// of the item within the heap. This is useful for the
// case where we want to replace an item in the heap and
// fix its structure; the container/heap.Fix method requires
// the index of the item that possibly violates the heap
// invariant.
type Item struct {
	Val   float64
	index int
}

// NewHeap initializes a new Heap.
func NewHeap(vals []float64, cmp func(float64, float64) bool) *Heap {
	items := []*Item{}
	for i, val := range vals {
		items = append(items, &Item{Val: val, index: i})
	}

	h := &Heap{items: items, cmp: cmp}
	heapops.Init(h)
	return h
}

func (h *Heap) Len() int {
	return len(h.items)
}

func (h *Heap) Less(i, j int) bool {
	return h.cmp(h.items[i].Val, h.items[j].Val)
}

func (h *Heap) Swap(i, j int) {
	h.items[i], h.items[j] = h.items[j], h.items[i]
}

// Push adds an element to the heap.
// This satisfies heapops.Interface.
func (h *Heap) Push(x interface{}) {
	h.items = append(h.items, x.(*Item))
}

// Pop removes element Len() - 1 from the heap.
// This satisfies heapops.Interface.
func (h *Heap) Pop() interface{} {
	x := h.items[len(h.items)-1]
	h.items = h.items[:len(h.items)-1]
	return x
}

// Peek returns the element at the top of the heap,
// without popping it.
func (h *Heap) Peek() float64 {
	return h.items[0].Val
}

// Values returns the values stored by the heap; this is
// mainly for debugging purposes.
func (h *Heap) Values() []float64 {
	vals := []float64{}
	for _, item := range h.items {
		vals = append(vals, item.Val)
	}
	return vals
}
