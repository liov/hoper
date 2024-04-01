package main

import "fmt"

/*
双马网络科技笔试答题：
反转从位置 m 到 n 的链表。请使用一趟扫描完成反转。

说明:
1 ≤m≤n≤ 链表长度。

示例:

输入: 1->2->3->4->5->NULL, m = 2, n = 4
输出: 1->4->3->2->5->NULL
*/
type node struct {
	Val  int
	Next *node
}

func reverse(head *node, m, n int) *node {
	if m == n {
		return head
	}
	var prev *node
	next := head
	idx := 0
	for next != nil {
		idx++
		if idx == m {
			l2head := reverseHelper(next, n-m+1)
			if m == 1 {
				return l2head
			} else {
				prev.Next = l2head
				return head
			}
		}
		prev = next
		next = next.Next
	}
	return nil
}

func reverseHelper(head *node, n int) *node {
	var newhead *node
	next := head
	var times int
	for next != nil {
		tmp := next.Next
		next.Next = newhead
		newhead = next
		next = tmp
		times++
		if times == n {
			head.Next = next
			return newhead
		}
	}
	return head
}

func main() {
	l1 := &node{
		Val:  1,
		Next: nil,
	}
	l2 := &node{
		Val:  2,
		Next: nil,
	}
	l1.Next = l2
	l3 := &node{
		Val:  3,
		Next: nil,
	}
	l2.Next = l3
	l4 := &node{
		Val:  4,
		Next: nil,
	}
	l3.Next = l4
	l5 := &node{
		Val:  5,
		Next: nil,
	}
	l4.Next = l5
	l := reverse(l1, 2, 4)
	for l != nil {
		fmt.Println(l)
		l = l.Next
	}

}
