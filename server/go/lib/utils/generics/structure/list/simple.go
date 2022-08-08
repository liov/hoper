package list

type SingleNode[T any] struct {
	Data T
	Next *SingleNode[T]
}

type SimpleList[T any] struct {
	Head, Tail *SingleNode[T]
	Size       int
}

func NewSimpleList[T any]() SimpleList[T] {
	l := SimpleList[T]{}
	l.Head = nil //head指向头部结点
	l.Tail = nil //tail指向尾部结点
	l.Size = 0
	return l
}

func (l *SimpleList[T]) Len() int {
	return l.Size
}

func (l *SimpleList[T]) First() T {
	if l.Size == 0 {
		panic("list is empty")
		return *new(T)
	}
	return l.Head.Data
}

func (l *SimpleList[T]) Pop() T {
	if l.Size == 0 {
		panic("list is empty")
		return *new(T)
	}

	p := l.Head
	l.Head = p.Next
	if l.Size == 1 {
		l.Tail = nil
	}
	l.Size--
	return p.Data
}

func (l *SimpleList[T]) Push(v T) {
	node := &SingleNode[T]{v, nil}
	if l.Size == 0 {
		l.Head = node
		l.Tail = node
		l.Size++
		return
	}
	l.Tail.Next = node
	l.Tail = node
	l.Size++
}
