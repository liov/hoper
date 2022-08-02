package list

import (
	"fmt"

	"github.com/actliboy/hoper/server/go/lib/utils/log"
)

//链表结点
type Node struct {
	data       interface{}
	next, prev *Node
}

//链表
type LinkList struct {
	head, tail *Node
	size       int
}

//新建空链表，即创建Node指针head，用来指向链表第一个结点，初始为空
func CreateLinkList() LinkList {
	l := LinkList{}
	l.head = nil //head指向头部结点
	l.tail = nil //tail指向尾部结点
	l.size = 0
	return l
}

//是否为空链表
func (l *LinkList) IsEmpty() bool {
	return l.size == 0
}

//获取链表长度
func (l *LinkList) GetLength() int {
	return l.size
}

//是否含有指定结点
func (l *LinkList) Exist(node *Node) bool {
	var p *Node = l.head
	for p != nil {
		if p == node {
			return true
		} else {
			p = p.next
		}
	}
	return false
}

//获取含有指定数据的第一个结点
func (l *LinkList) GetNode(e interface{}) *Node {
	var p *Node = l.head
	for p != nil {
		//找到该数据所在结点
		if e == p.data {
			return p
		} else {
			p = p.next
		}
	}
	return nil
}

//在链表尾部添加数据
func (l *LinkList) Append(e interface{}) {
	//为数据创建新结点
	newNode := Node{}
	newNode.data = e
	newNode.next = nil

	if l.size == 0 {
		l.head = &newNode
		l.tail = &newNode
	} else {
		l.tail.next = &newNode
		l.tail = &newNode
	}
	l.size++
}

//在链表头部插入数据
func (l *LinkList) InsertHead(e interface{}) {
	newNode := Node{}
	newNode.data = e
	newNode.next = l.head
	l.head = &newNode
	if l.size == 0 {
		l.tail = &newNode
	}
	l.size++
}

//在指定结点后面插入数据
func (l *LinkList) InsertAfterNode(pre *Node, e interface{}) {
	//如果链表中存在该结点，才进行插入
	if l.Exist(pre) {
		newNode := Node{}
		newNode.data = e
		if pre.next == nil {
			l.Append(e)
		} else {
			newNode.next = pre.next
			pre.next = &newNode
		}
		l.size++
	} else {
		log.Error("链表中不存在该结点")
	}
}

//在第一次出现指定数据的结点后插入数据,若链表中无该数据，返回false
func (l *LinkList) InsertAfterData(preData interface{}, e interface{}) bool {
	var p *Node = l.head
	for p != nil {
		//找到该数据所在结点
		if p.data == preData {
			l.InsertAfterNode(p, e)
			return true
		} else {
			p = p.next
		}
	}
	//没有找到该数据
	log.Error("链表中没有该数据，插入失败")
	return false
}

//在指定下标处插入数据
func (l *LinkList) Insert(position int, e interface{}) bool {
	if position < 0 {
		log.Error("指定下标不合法")
		return false
	} else if position == 0 {
		//在头部插入
		l.InsertHead(e)
		return true
	} else if position == l.size {
		//在尾部插入
		l.Append(e)
		return true
	} else if position > l.size {
		log.Error("指定下标超出链表长度")
		return false
	} else {
		//在中间插入
		var index int
		var p *Node = l.head
		//逐个移动指针
		//position是插入后新结点的下标，position-1时需要定位到的其前一个结点的下标
		for index = 0; index < position-1; index++ {
			p = p.next
		}
		//找到
		l.InsertAfterNode(p, e)
		return true
	}

}

//删除指定结点
func (l *LinkList) DeleteNode(node *Node) {
	//存在该结点
	if l.Exist(node) {
		//如果是头部结点
		if node == l.head {
			l.head = l.head.next
			//如果是尾部结点
		} else if node == l.tail {
			//寻找指向其前一个结点的指针
			var p *Node = l.head
			for p.next != l.tail {
				p = p.next
			}
			p.next = nil
			l.tail = p
			//中间结点
		} else {
			var p *Node = l.head
			for p.next != node {
				p = p.next
			}
			p.next = node.next
		}
		l.size--
	}
}

//删除第一个含指定数据的结点
func (l *LinkList) Delete(e interface{}) {
	p := l.GetNode(e)
	if p == nil {
		log.Error("链表中无该数据，删除失败")
	} else {
		l.DeleteNode(p)
	}
}

//遍历链表
func (l *LinkList) Traverse(f func(interface{})) {
	var p *Node = l.head
	if l.IsEmpty() {
		log.Error("LinkList is empty")
	} else {
		for p != nil {
			if p.data != nil {
				f(p.data)
			}
			p = p.next
		}
	}
}

//打印链表信息
func (l *LinkList) PrintInfo() {
	fmt.Println("###############################################")
	fmt.Println("链表长度为：", l.GetLength())
	fmt.Println("链表是否为空:", l.IsEmpty())
	fmt.Print("遍历链表：")
	l.Traverse(func(data interface{}) { log.Info(data, " ") })
	fmt.Println("###############################################")
}
