package aop

func Aop(before, fs, after []func()) {
	for _, f := range before {
		f()
	}
	for _, f := range fs {
		f()
	}
	for _, f := range after {
		f()
	}
}
