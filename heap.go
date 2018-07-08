package stream

import heapops "container/heap"

type heap struct {
	max  bool
	vals []float64
}

func newHeap(vals []float64, max bool) *heap {
	h := &heap{max: max, vals: vals}
	heapops.Init(h)
	return h
}

func (h *heap) Len() int {
	return len(h.vals)
}

func (h *heap) Less(i, j int) bool {
	if h.max {
		return h.vals[i] > h.vals[j]
	}

	return h.vals[i] < h.vals[j]
}

func (h *heap) Swap(i, j int) {
	h.vals[i], h.vals[j] = h.vals[j], h.vals[i]
}

func (h *heap) Push(x interface{}) {
	h.vals = append(h.vals, x.(float64))
}

func (h *heap) Pop() interface{} {
	x := h.vals[len(h.vals)-1]
	h.vals = h.vals[:len(h.vals)-1]
	return x
}

func (h *heap) peek() float64 {
	return h.vals[0]
}
