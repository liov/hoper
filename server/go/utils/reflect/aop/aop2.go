package aop

type AnyFunc func()

func (f AnyFunc) Aop(before, after AnyFunc) AnyFunc {
	return func() {
		before()
		f()
		after()
	}
}
