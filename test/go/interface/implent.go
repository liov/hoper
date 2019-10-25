package main

type FooInter interface {
	Like(inter FooInter)
}

type FooSlice []FooInter

func (f *FooSlice) Like(inter FooInter)  {
	f.Like(inter)
}