package lazy

import _interface "github.com/liov/hoper/server/go/lib/utils/def/interface"

type Lazy[T _interface.Init] struct {
	init bool
	Prop T
}

func (l *Lazy[T]) GetProp() T {
	if l.init {
		return l.Prop
	}
	l.Prop.Init()
	l.init = true
	return l.Prop
}
