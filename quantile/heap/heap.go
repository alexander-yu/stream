package heap

import heapops "container/heap"

// Heap implements a heap data structure.
type Heap struct {
	vals []interface{}
	cmp  func(interface{}, interface{}) bool
}

// NewHeap initializes a new Heap.
func NewHeap(vals []interface{}, cmp func(interface{}, interface{}) bool) *Heap {
	h := &Heap{vals: vals, cmp: cmp}
	heapops.Init(h)
	return h
}

func (h *Heap) Len() int {
	return len(h.vals)
}

func (h *Heap) Less(i, j int) bool {
	return h.cmp(h.vals[i], h.vals[j])
}

func (h *Heap) Swap(i, j int) {
	h.vals[i], h.vals[j] = h.vals[j], h.vals[i]
}

// Push adds an element to the heap.
// This satisfies heapops.Interface.
func (h *Heap) Push(x interface{}) {
	h.vals = append(h.vals, x)
}

// Pop removes element Len() - 1 from the heap.
// This satisfies heapops.Interface.
func (h *Heap) Pop() interface{} {
	x := h.vals[len(h.vals)-1]
	h.vals = h.vals[:len(h.vals)-1]
	return x
}

// Peek returns the element at the top of the heap,
// without popping it.
func (h *Heap) Peek() interface{} {
	return h.vals[0]
}

// Values returns the values stored by the heap; this is
// mainly for debugging purposes.
func (h *Heap) Values() []interface{} {
	return h.vals
}
