package main

func main() {
	//GetFoo() = nil Cannot assign to GetFoo()
	/*Unused variable 'f',必须用一下，例如fmt.Println(f)
	f:=GetFoo()
	f = nil*/
}

type Foo struct {
}

func GetFoo() *Foo {
	return &Foo{}
}
