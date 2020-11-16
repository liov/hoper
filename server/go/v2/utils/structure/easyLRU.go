package structure

// Node ...
type Node struct {
	k, v int
}

// EasyLRU ...
type EasyLRU struct {
	size  int
	store []*Node
}

// NewLRU ...
func NewLRU(capacity int) *EasyLRU {
	return &EasyLRU{
		size:  capacity,
		store: make([]*Node, capacity),
	}
}

// Get ...
func (l *EasyLRU) Get(key int) int {
	if key < 0 {
		return -1
	}
	for i := range l.store {
		if l.store[i] == nil {
			break
		}
		if key == l.store[i].k {
			n := l.store[i]
			copy(l.store[1:i+1], l.store[0:i])
			l.store[0] = n
			return n.v
		}
	}
	return -1
}

// Put ...
func (l *EasyLRU) Put(key int, value int) {
	// lookup
	for i := range l.store {
		if l.store[i] == nil {
			break
		}
		// found
		if key == l.store[i].k {
			n := l.store[i]
			copy(l.store[1:i+1], l.store[0:i])
			n.k = key
			n.v = value
			l.store[0] = n
			return
		}
	}

	// not found
	if l.size == len(l.store) {
		copy(l.store[1:l.size], l.store[0:l.size-1])
	} else {
		copy(l.store[1:l.size+1], l.store[0:l.size])
		l.size++
	}
	l.store[0] = &Node{k: key, v: value}
}

