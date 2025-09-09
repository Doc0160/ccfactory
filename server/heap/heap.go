package heap

// Heap is a binary heap. Comparator returns true if a < b (i.e. a has higher priority).
type Heap[T any] struct {
	data []T
	less func(a, b T) bool
}

// New creates an empty heap. less(a,b)==true means a is higher priority than b (min-heap if less = a<b).
func New[T any](less func(a, b T) bool) *Heap[T] {
	return &Heap[T]{less: less}
}

// FromSlice builds a heap from a slice (O(n)).
func FromSlice[T any](items []T, less func(a, b T) bool) *Heap[T] {
	h := &Heap[T]{data: append([]T(nil), items...), less: less}
	n := len(h.data)
	for i := parent(n - 1); i >= 0; i-- {
		h.down(i, n)
	}
	return h
}

func parent(i int) int { return (i - 1) / 2 }
func left(i int) int   { return 2*i + 1 }
func right(i int) int  { return 2*i + 2 }

// Len returns number of elements.
func (h *Heap[T]) Len() int { return len(h.data) }

// IsEmpty returns whether heap is empty.
func (h *Heap[T]) IsEmpty() bool { return len(h.data) == 0 }

// Peek returns top element without removing. ok==false if empty.
func (h *Heap[T]) Peek() (top T, ok bool) {
	if h.IsEmpty() {
		var zero T
		return zero, false
	}
	return h.data[0], true
}

// Push adds an element.
func (h *Heap[T]) Push(x T) {
	h.data = append(h.data, x)
	h.up(len(h.data) - 1)
}

// Pop removes and returns top element. ok==false if empty.
func (h *Heap[T]) Pop() (top T, ok bool) {
	if h.IsEmpty() {
		var zero T
		return zero, false
	}
	n := len(h.data)
	h.swap(0, n-1)
	top = h.data[n-1]
	h.data = h.data[:n-1]
	if len(h.data) > 0 {
		h.down(0, len(h.data))
	}
	return top, true
}

// Replace pops and pushes atomically: returns old top (ok=false if empty).
func (h *Heap[T]) Replace(x T) (old T, ok bool) {
	if h.IsEmpty() {
		h.Push(x)
		var zero T
		return zero, false
	}
	old = h.data[0]
	h.data[0] = x
	h.down(0, len(h.data))
	return old, true
}

func (h *Heap[T]) up(i int) {
	for i > 0 {
		p := parent(i)
		if !h.less(h.data[i], h.data[p]) {
			break
		}
		h.swap(i, p)
		i = p
	}
}

func (h *Heap[T]) down(i, n int) bool {
	moved := false
	for {
		l := left(i)
		if l >= n {
			return moved
		}
		smallest := l
		r := l + 1
		if r < n && h.less(h.data[r], h.data[l]) {
			smallest = r
		}
		if !h.less(h.data[smallest], h.data[i]) {
			return moved
		}
		h.swap(i, smallest)
		i = smallest
		moved = true
	}
}

func (h *Heap[T]) Fix(i int) {
	if !h.down(i, len(h.data)) {
		h.up(i)
	}
}

func (h *Heap[T]) swap(i, j int) { h.data[i], h.data[j] = h.data[j], h.data[i] }

/* -------------------------
   Example usage (not part of package):
---------------------------------

package main

import (
	"fmt"
	"your/module/heap"
)

func main() {
	// min-heap for ints
	h := heap.New(func(a, b int) bool { return a < b })
	h.Push(5); h.Push(2); h.Push(9)
	fmt.Println(h.Pop()) // 2,true
	fmt.Println(h.Peek()) // 5,true

	// max-heap for ints
	max := heap.New(func(a, b int) bool { return a > b })
	for _, v := range []int{3,1,4,1,5} { max.Push(v) }
	for !max.IsEmpty() {
		v, _ := max.Pop()
		fmt.Println(v) // 5,4,3,1,1
	}
}

------------------------- */
