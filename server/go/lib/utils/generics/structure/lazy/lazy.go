package lazy

type InitInterface interface {
	Init()
}

type Lazy[T InitInterface] struct {
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
