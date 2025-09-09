package factory

import (
	"container/heap"
	"encoding/json"
)

type Less interface {
	Less(other Less) bool
}

type binaryHeapInternal[T Less] struct {
	data []T
}

func (h *binaryHeapInternal[T]) Len() int {
	return len(h.data)
}

func (h *binaryHeapInternal[T]) Less(i, j int) bool {
	return h.data[i].Less(h.data[j])
}

func (h *binaryHeapInternal[T]) Swap(i, j int) {
	h.data[i], h.data[j] = h.data[j], h.data[i]
}

func (h *binaryHeapInternal[T]) Push(x any) {
	h.data = append(h.data, x.(T))
}
func (h *binaryHeapInternal[T]) Pop() any {
	old := h.data
	n := len(old)
	x := old[n-1]
	h.data = old[:n-1]
	return x
}

type BinaryHeap[T Less] struct {
	data binaryHeapInternal[T]
}

func NewBinaryHeap[T Less]() BinaryHeap[T] {
	h := BinaryHeap[T]{
		data: binaryHeapInternal[T]{
			data: []T{},
		},
	}
	heap.Init(&h.data)
	return h
}

func (h BinaryHeap[T]) Push(x T) {
	heap.Push(&h.data, x)
}

func (h *BinaryHeap[T]) Remove() T {
	return heap.Pop(&h.data).(T)
}

type RawMessage = json.RawMessage

func Into[T any](r RawMessage) (T, error) {
	var v T
	err := json.Unmarshal(r, &v)
	return v, err
}
