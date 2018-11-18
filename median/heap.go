package median

import heapops "container/heap"

type heap struct {
	vals []interface{}
	cmp  func(interface{}, interface{}) bool
}

func newHeap(vals []interface{}, cmp func(interface{}, interface{}) bool) *heap {
	h := &heap{vals: vals, cmp: cmp}
	heapops.Init(h)
	return h
}

func (h *heap) Len() int {
	return len(h.vals)
}

func (h *heap) Less(i, j int) bool {
	return h.cmp(h.vals[i], h.vals[j])
}

func (h *heap) Swap(i, j int) {
	h.vals[i], h.vals[j] = h.vals[j], h.vals[i]
}

func (h *heap) Push(x interface{}) {
	h.vals = append(h.vals, x)
}

func (h *heap) Pop() interface{} {
	x := h.vals[len(h.vals)-1]
	h.vals = h.vals[:len(h.vals)-1]
	return x
}

func (h *heap) peek() interface{} {
	return h.vals[0]
}
