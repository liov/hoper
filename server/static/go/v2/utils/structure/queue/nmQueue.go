package queue

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

type NmQueue struct {
	head *Node
	end  *Node
	size int
}

func NewNmQueue() *NmQueue {
	q := &NmQueue{nil, nil, nil}
	return q
}

func (q *NmQueue) push(data interface{}) {
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

func (q *NmQueue) pop() (interface{}, bool) {
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
