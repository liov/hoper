package queue

type Node struct {
	data     interface{}
	previous *Node
	next     *Node
}

func (n *Node) Value() interface{} {
	return n.data
}

func (n *Node) Set(value interface{}) {
	n.data = value
}

func (n *Node) Previous() *Node {
	return n.previous
}

func (n *Node) Next() *Node {
	return n.next
}

type ListQueue struct {
	head *Node
	end  *Node
	size int
}

func NewListQueue(size int) *ListQueue {
	q := &ListQueue{nil, nil, size}
	return q
}

func (q *ListQueue) push(data interface{}) {
	n := &Node{data: data, next: nil}

	if q.end == nil {
		q.head = n
		q.end = n
	} else {
		q.end.next = n
		q.end = n
	}

	return
}

func (q *ListQueue) pop() (interface{}, bool) {
	if q.head == nil {
		return nil, false
	}

	data := q.head.data
	q.head = q.head.next
	if q.head == nil {
		q.end = nil
	}

	return data, true
}
