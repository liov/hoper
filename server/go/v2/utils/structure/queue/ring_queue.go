package queue

type Queue interface {
	// 获取当前链表长度。
	Len() int
	// 获取当前链表容量。
	Capacity() int
	// 获取当前链表头结点。
	Front() *Node
	// 获取当前链表尾结点。
	Rear() *Node
	// 入列。
	Enqueue(value interface{}) bool
	// 出列。
	Dequeue() interface{}
}

type RingQueue struct {
	head, tail int
	len        int
	buf        []interface{}
}

func NewRingQueue(capacity int) *RingQueue {
	nodes := make([]interface{}, capacity, capacity)
	return &RingQueue{
		head: -1,
		tail: -1,
		buf:  nodes,
	}
}

func (q *RingQueue) Length() int {
	return q.len
}

func (q *RingQueue) Capacity() int {
	return len(q.buf)
}

func (q *RingQueue) Front() interface{} {
	if q.len == 0 {
		return nil
	}

	return q.buf[q.head]
}

func (q *RingQueue) Tail() interface{} {
	if q.len == 0 {
		return nil
	}

	return q.buf[q.tail]
}

func (q *RingQueue) Enqueue(value interface{}) bool {
	if q.IsFull() || value == nil {
		return false
	}

	q.tail++
	if q.tail == len(q.buf) {
		q.tail = 0
	}
	q.buf[q.tail] = value
	q.len++

	if q.len == 1 {
		q.head = q.tail
	}

	return true
}

func (q *RingQueue) Dequeue() interface{} {
	if q.len == 0 {
		return nil
	}

	result := q.buf[q.head]
	q.buf[q.head] = nil
	q.head++
	q.len--
	if q.head == len(q.buf) {
		q.head = 0
	}

	return result
}


// IsFull checks if the ring buffer is full
func (q *RingQueue) IsFull() bool {
	return q.len == len(q.buf)
}


// LookAll reads all elements from ring buffer
// this method doesn't consume all elements
func (q *RingQueue) LookAll() []interface{} {
	all := make([]interface{}, q.len)
	if q.len == 0{
		return all
	}
	j := 0
	for i := q.head; ; i++ {
		if i == len(q.buf) {
			i = 0
		}
		all[j] = q.buf[i]
		if i == q.tail {
			break
		}
		j++
	}
	return all
}
