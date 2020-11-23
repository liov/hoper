package pool

// chan pool 固定大小
type Pool struct {
	using []interface{}
	pool  []interface{}
}

// 链表池，用放入列表尾，用完放入队列头，去头的时候判断，如果都是使用中，新建加入队头
type Pool2 struct {
}
