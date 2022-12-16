package heap

import "testing"

type Foo struct {
	A int
}

func (f *Foo) CompareField() int {
	return f.A
}

func TestHeap(t *testing.T) {
	heap := Heap[*Foo]{}
	heap.Init()
	heap.Push(&Foo{10})
	heap.Push(&Foo{5})
	heap.Push(&Foo{8})
	heap.Push(&Foo{2})
	heap.Push(&Foo{26})
	heap.Push(&Foo{6})
	heap.Push(&Foo{9})
	heap.Push(&Foo{1})
	for _, foo := range heap {
		t.Log(foo)
	}
}
