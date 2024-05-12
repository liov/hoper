package main

type Node struct {
	Key, Val   int
	Prev, Next *Node
}
type LRUCache struct {
	capacity   int
	data       map[int]*Node
	Head, Tail *Node
}

func New(capacity int) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		data:     make(map[int]*Node),
		Head:     nil,
		Tail:     nil,
	}
}

func (l *LRUCache) Get(key int) int {
	v, ok := l.data[key]
	if ok {
		l.insert(v)
		return v.Val
	}
	return 0
}

func (l *LRUCache) insert(v *Node) {
	v.Prev.Next = v.Next
	v.Next.Prev = v.Prev
	head := l.Head
	head.Prev = v
	l.Head = v
	v.Next = head
}

func (l *LRUCache) Put(key, val int) {
	v, ok := l.data[key]
	if ok {
		l.insert(v)
		v.Val = val
	}

	if len(l.data) == l.capacity {
		delete(l.data, l.Tail.Key)
		l.Tail = l.Tail.Prev
	}

	newNode := &Node{Key: key, Val: val}
	newNode.Next = l.Head
	l.Head.Prev = newNode
	l.Head = newNode
	l.data[key] = newNode
}
