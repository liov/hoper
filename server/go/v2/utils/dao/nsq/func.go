package tnsq

type Foo struct {
}

type Bar struct {
}

func (b *Bar) BarFuncAdd(argOne, argTwo float64) float64 {

	return argOne + argTwo
}

func (f *Foo) FooFuncSwap(argOne, argTwo string) (string, string) {

	return argTwo, argOne
}
