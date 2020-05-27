package queue

import (
	"container/list"
	"errors"
)

type Queue interface {
	// 获取当前链表长度。
	Length() int
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

type Node struct {
	data     interface{}
	previous *Node
	next     *Node
}

type MyList list.List

func (m *MyList) Lenght() int {
	return list.List(*m).Len()
}

type MyQueue struct {
	front    int
	rear     int
	length   int
	capacity int
	nodes    []*Node
}

func NewMyQueue(capacity int) (*MyQueue, error) {
	if capacity <= 0 {
		return nil, errors.New("capacity is less than 0")
	}

	nodes := make([]*Node, capacity, capacity)

	return &MyQueue{
		front:    -1,
		rear:     -1,
		capacity: capacity,
		nodes:    nodes,
	}, nil
}

func (q *MyQueue) Length() int {
	return q.length
}

func (q *MyQueue) Capacity() int {
	return q.capacity
}

func (q *MyQueue) Front() *Node {
	if q.length == 0 {
		return nil
	}

	return q.nodes[q.front]
}

func (q *MyQueue) Rear() *Node {
	if q.length == 0 {
		return nil
	}

	return q.nodes[q.rear]
}

func (q *MyQueue) Enqueue(value interface{}) bool {
	if q.length == q.capacity || value == nil {
		return false
	}

	node := &Node{
		data: value,
	}

	index := (q.rear + 1) % cap(q.nodes)
	q.nodes[index] = node
	q.rear = index
	q.length++

	if q.length == 1 {
		q.front = index
	}

	return true
}

func (q *MyQueue) Dequeue() interface{} {
	if q.length == 0 {
		return nil
	}

	result := q.nodes[q.front].data
	q.nodes[q.front] = nil
	index := (q.front + 1) % cap(q.nodes)
	q.front = index
	q.length--

	return result
}
