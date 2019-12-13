package aop

func Aop(fs []func()) {
	for _, f := range fs {
		f()
	}
}
