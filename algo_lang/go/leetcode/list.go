package leetcode

import "strconv"

type ListNode struct {
	Val  int
	Next *ListNode
}

func NewList(arr []int) *ListNode {
	head := &ListNode{Val: arr[0]}
	tmp := head
	for i := 1; i < len(arr); i++ {
		tmp.Next = &ListNode{Val: arr[i]}
		tmp = tmp.Next
	}
	return head
}

func (list *ListNode) String() string {
	var s string
	l := list
	for l != nil {
		s += strconv.Itoa(l.Val) + ","
		l = l.Next
	}
	return s
}

// Node Val > 0
type Node struct {
	Val   int
	Prev  *Node
	Next  *Node
	Child *Node
}

func NewNode(arr []int) *Node {
	return nil
}
