package list

type Node[T any] struct {
	Data T
	Next *Node[T]
}

type List[T any] struct {
	Head, Tail *Node[T]
	Size       uint
}

func New[T any]() List[T] {
	l := List[T]{}
	l.Head = nil //head指向头部结点
	l.Tail = nil //tail指向尾部结点
	l.Size = 0
	return l
}

func (l *List[T]) Len() uint {
	return l.Size
}

func (l *List[T]) First() T {
	if l.Size == 0 {
		panic("list is empty")
		return *new(T)
	}
	return l.Head.Data
}

func (l *List[T]) Pop() T {
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

func (l *List[T]) Push(v T) {
	node := &Node[T]{v, nil}
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
