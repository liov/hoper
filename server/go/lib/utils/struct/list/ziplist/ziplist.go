package ziplist

import (
	"encoding/binary"
	"errors"
	"sync"
	"unsafe"
)

// TODO: fix

const (
	ZIPLIST_ENCODING_RAW = 0x00
	ZIPLIST_ENCODING_INT = 0x01
)

type zlentry struct {
	// 编码方式，可以是ZIPLIST_ENCODING_RAW或ZIPLIST_ENCODING_INT
	encoding byte
	// 数据指针，指向数据所在的内存地址
	ptr unsafe.Pointer
	// 数据长度，如果是int类型则为8，否则为实际数据长度
	length     uint32
	next, prev *zlentry
}

type ziplist struct {
	// 内存池，用于分配和释放内存
	pool sync.Pool
	// 元素数量
	length uint32
	// 尾节点指针
	tail *zlentry
	// 头节点指针
	head *zlentry
}

// ZlValue表示ziplist中节点的值
type ZlValue interface{}

// ZlEntry表示ziplist中的一个节点
type ZlEntry struct {
	Value ZlValue
}

// 用于存储int类型的值
type intVal struct {
	Value int64
}

// 用于存储字符串类型的值
type stringVal struct {
	Value []byte
}

// ziplist节点中存储int类型值的最大长度
const maxIntValLen = 9

// ziplist节点中存储字符串类型值的最大长度
const maxStringValLen = 255

// 新建一个ziplist
func newZiplist() *ziplist {
	z := &ziplist{}
	z.pool.New = func() interface{} {
		return make([]byte, 1024)
	}
	return z
}

// 在ziplist中添加一个元素
func (z *ziplist) push(value []byte) error {
	// 检查value长度是否超过限制
	if len(value) > 255 {
		return errors.New("value too long")
	}

	// 计算新节点的长度
	nodeLength := uint32(len(value) + 1)
	if nodeLength < 254 {
		nodeLength++
	}

	// 分配内存
	buf := z.pool.Get().([]byte)[:nodeLength]

	// 写入数据
	buf[0] = byte(len(value))
	copy(buf[1:], value)

	// 新建一个节点
	node := &zlentry{
		encoding: ZIPLIST_ENCODING_RAW,
		ptr:      unsafe.Pointer(&buf[0]),
		length:   nodeLength,
	}

	// 更新tail指针
	if z.tail != nil {
		z.tail = node
	}

	// 更新head指针
	if z.head == nil {
		z.head = node
	}

	// 更新元素数量
	z.length++

	return nil
}

// 在ziplist中添加一个int类型的元素
func (z *ziplist) pushInt(value int64) error {
	// 计算新节点的长度
	nodeLength := uint32(9)

	// 分配内存
	buf := z.pool.Get().([]byte)[:nodeLength]

	// 写入数据
	buf[0] = 0xff
	binary.LittleEndian.PutUint64(buf[1:], uint64(value))

	// 新建一个节点
	node := &zlentry{
		encoding: ZIPLIST_ENCODING_INT,
		ptr:      unsafe.Pointer(&buf[0]),
		length:   nodeLength,
	}

	// 更新tail指针
	if z.tail != nil {
		z.tail = node
	}

	// 更新head指针
	if z.head == nil {
		z.head = node
	}

	// 更新元素数量
	z.length++

	return nil
}

// 压缩ziplist，删除

// 删除节点，如果节点是raw类型，则直接释放内存；
// 如果是int类型，则不需要释放内存
func (z *ziplist) deleteNode(node *zlentry) {
	if node.encoding == ZIPLIST_ENCODING_RAW {
		buf := (*[1 << 30]byte)(node.ptr)[:node.length]
		z.pool.Put(buf)
	}

	if node == z.head {
		z.head = node.next
	}

	if node == z.tail {
		z.tail = node.prev
	}

	z.length--
}

// 在指定节点之后插入新节点
func (z *ziplist) insertAfter(node *zlentry, value []byte) error {
	if node == nil {
		return errors.New("node is nil")
	}

	// 检查value长度是否超过限制
	if len(value) > 255 {
		return errors.New("value too long")
	}

	// 计算新节点的长度
	nodeLength := uint32(len(value) + 1)
	if nodeLength < 254 {
		nodeLength++
	}

	// 分配内存
	buf := z.pool.Get().([]byte)[:nodeLength]

	// 写入数据
	buf[0] = byte(len(value))
	copy(buf[1:], value)

	// 新建一个节点
	newNode := &zlentry{
		encoding: ZIPLIST_ENCODING_RAW,
		ptr:      unsafe.Pointer(&buf[0]),
		length:   nodeLength,
	}

	// 插入节点
	newNode.prev = node
	node.next = newNode

	// 更新tail指针
	if node == z.tail {
		z.tail = newNode
	}

	// 更新元素数量
	z.length++

	return nil
}

// 在指定节点之前插入新节点
func (z *ziplist) insertBefore(node *zlentry, value []byte) error {
	if node == nil {
		return errors.New("node is nil")
	}

	// 检查value长度是否超过限制
	if len(value) > 255 {
		return errors.New("value too long")
	}

	// 计算新节点的长度
	nodeLength := uint32(len(value) + 1)
	if nodeLength < 254 {
		nodeLength++
	}

	// 分配内存
	buf := z.pool.Get().([]byte)[:nodeLength]

	// 写入数据
	buf[0] = byte(len(value))
	copy(buf[1:], value)

	// 新建一个节点
	newNode := &zlentry{
		encoding: ZIPLIST_ENCODING_RAW,
		ptr:      unsafe.Pointer(&buf[0]),
		length:   nodeLength,
	}

	// 查找node的前置节点
	prev := z.head
	for prev != nil && prev.next != node {
		prev = prev.next
	}

	// 插入节点
	newNode.next = node
	prev.next = newNode

	// 更新head指针
	if node == z.head {
		z.head = newNode
	}

	// 更新元素数量
	z.length++

	return nil
}
